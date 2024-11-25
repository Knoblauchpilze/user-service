package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/backend-toolkit/pkg/db"
	"github.com/KnoblauchPilze/backend-toolkit/pkg/errors"
	"github.com/KnoblauchPilze/user-service/pkg/communication"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, userDto communication.UserDtoRequest) (communication.UserDtoResponse, error)
	Get(ctx context.Context, id uuid.UUID) (communication.UserDtoResponse, error)
	List(ctx context.Context) ([]uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, userDto communication.UserDtoRequest) (communication.UserDtoResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Login(ctx context.Context, userDto communication.UserDtoRequest) (communication.ApiKeyDtoResponse, error)
	Logout(ctx context.Context, id uuid.UUID) error
}

type userServiceImpl struct {
	conn db.Connection

	userRepo   repositories.UserRepository
	apiKeyRepo repositories.ApiKeyRepository

	apiKeyValidity time.Duration
}

func NewUserService(config ApiKeyConfig, conn db.Connection, repos repositories.Repositories) UserService {
	return &userServiceImpl{
		conn:       conn,
		userRepo:   repos.User,
		apiKeyRepo: repos.ApiKey,

		apiKeyValidity: config.Validity,
	}
}

func (s *userServiceImpl) Create(ctx context.Context, userDto communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	user := communication.FromUserDtoRequest(userDto)

	if user.Email == "" {
		return communication.UserDtoResponse{}, errors.NewCode(InvalidEmail)
	}
	if user.Password == "" {
		return communication.UserDtoResponse{}, errors.NewCode(InvalidPassword)
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	out := communication.ToUserDtoResponse(createdUser)
	return out, nil
}

func (s *userServiceImpl) Get(ctx context.Context, id uuid.UUID) (communication.UserDtoResponse, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	out := communication.ToUserDtoResponse(user)
	return out, nil
}

func (s *userServiceImpl) List(ctx context.Context) ([]uuid.UUID, error) {
	return s.userRepo.List(ctx)
}

func (s *userServiceImpl) Update(ctx context.Context, id uuid.UUID, userDto communication.UserDtoRequest) (communication.UserDtoResponse, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	user.Email = userDto.Email
	user.Password = userDto.Password

	updated, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return communication.UserDtoResponse{}, err
	}

	out := communication.ToUserDtoResponse(updated)
	return out, nil
}

func (s *userServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	err = s.apiKeyRepo.DeleteForUser(ctx, tx, id)
	if err != nil {
		return err
	}
	err = s.userRepo.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) Login(ctx context.Context, user communication.UserDtoRequest) (communication.ApiKeyDtoResponse, error) {
	dbUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return communication.ApiKeyDtoResponse{}, err
	}

	if user.Password != dbUser.Password {
		return communication.ApiKeyDtoResponse{}, errors.NewCode(InvalidCredentials)
	}

	apiKey := persistence.ApiKey{
		Id:         uuid.New(),
		Key:        uuid.New(),
		ApiUser:    dbUser.Id,
		ValidUntil: time.Now().Add(s.apiKeyValidity),
	}

	createdKey, err := s.apiKeyRepo.Create(ctx, apiKey)
	if err != nil {
		return communication.ApiKeyDtoResponse{}, err
	}

	out := communication.ToApiKeyDtoResponse(createdKey)
	return out, nil
}

func (s *userServiceImpl) Logout(ctx context.Context, id uuid.UUID) error {
	_, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	tx, err := s.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Close(ctx)

	return s.apiKeyRepo.DeleteForUser(ctx, tx, id)
}
