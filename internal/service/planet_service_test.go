package service

import (
	"context"
	"testing"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/communication"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var defaultPlanetId = uuid.MustParse("5b0efd85-8817-4454-b8f3-7af5d93253a1")
var defaultPlanetName = "my-planet"
var defaultBuildingActionId = uuid.MustParse("38a739bd-79db-453e-ab03-44f9f300c3c8")

var defaultPlanetDtoRequest = communication.PlanetDtoRequest{
	Player: defaultPlayerId,
	Name:   defaultPlanetName,
}
var defaultPlanet = persistence.Planet{
	Id:     defaultPlanetId,
	Player: defaultPlayerId,
	Name:   defaultPlanetName,

	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultPlanetResource = persistence.PlanetResource{
	Planet:    defaultPlanetId,
	Resource:  metalResourceId,
	Amount:    784.0987,
	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultPlanetResourceProduction = persistence.PlanetResourceProduction{
	Planet:     defaultPlanetId,
	Building:   &defaultBuildingId,
	Resource:   metalResourceId,
	Production: 31,
	CreatedAt:  testDate,
	UpdatedAt:  testDate,
}
var defaultPlanetResourceStorage = persistence.PlanetResourceStorage{
	Planet:    defaultPlanetId,
	Resource:  metalResourceId,
	Storage:   9876,
	CreatedAt: testDate,
	UpdatedAt: testDate,
}
var defaultPlanetBuilding = persistence.PlanetBuilding{
	Planet:    defaultPlanetId,
	Building:  defaultBuildingId,
	Level:     2,
	CreatedAt: testDate,
	UpdatedAt: testDate,
}

func Test_PlanetService(t *testing.T) {
	s := ServicePoolTestSuite{
		generateRepositoriesMocks:      generatePlanetServiceMocks,
		generateErrorRepositoriesMocks: generateErrorPlanetServiceMocks,

		repositoryInteractionTestCases: map[string]repositoryInteractionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Create(ctx, defaultPlanetDtoRequest)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.createCalled)
					assert.Equal(defaultPlanetDtoRequest.Player, m.createdPlanet.Player)
					assert.Equal(defaultPlanetDtoRequest.Name, m.createdPlanet.Name)
				},
			},
			"create_repositoryFails": {
				generateRepositoriesMocks: generateErrorPlanetServiceMocks,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Create(ctx, defaultPlanetDtoRequest)
					return err
				},
				expectedError: errDefault,
			},
			"get_planet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.getCalled)
					assert.Equal(defaultPlanetId, m.getId)
				},
			},
			"get_planetRepositoryFails": {
				generateRepositoriesMocks: generateErrorPlanetServiceMocks,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_planetResource": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultPlanetId}, m.listForPlanetIds)
				},
			},
			"get_planetResourceRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_planetResourceProduction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultPlanetId}, m.listForPlanetIds)
				},
			},
			"get_planetResourceProductionRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_planetResourceStorage": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceStorageRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal([]uuid.UUID{defaultPlanetId}, m.listForPlanetIds)
				},
			},
			"get_planetResourceStorageRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_planetBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultPlanetId, m.listForPlanetId)
				},
			},
			"get_planetBuildingRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"get_buildingAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlanetCalled)
					assert.Equal(defaultPlanetId, m.listForPlanetId)
				},
			},
			"get_buildingActionRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
				expectedError: errDefault,
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.listCalled)
				},
			},
			"list_planetRepositoryFails": {
				generateRepositoriesMocks: generateErrorPlanetServiceMocks,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
				expectedError: errDefault,
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.listForPlayerCalled)
					assert.Equal(defaultPlayerId, m.listForPlayerId)
				},
			},
			"listForPlayer_planetRepositoryFails": {
				generateRepositoriesMocks: generateErrorPlanetServiceMocks,
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
				expectedError: errDefault,
			},
			"delete_buildingActionResourceProduction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_buildingActionCost": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_buildingAction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_planetBuilding": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_planetResourceProduction": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_planetResourceStorage": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceStorageRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_planetResource": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
					assert.Equal(defaultPlanetId, m.deleteForPlanetId)
				},
			},
			"delete_planet": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
					assert.Equal(defaultPlanetId, m.deleteId)
				},
			},
			"delete_buildingActionResourceProductionRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.BuildingActionResourceProduction = &mockBuildingActionResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
				},
			},
			"delete_buildingActionCostRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.BuildingActionCost = &mockBuildingActionCostRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionCostRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
				},
			},
			"delete_buildingActionRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.BuildingAction = &mockBuildingActionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertBuildingActionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
				},
			},
			"delete_planetBuildingRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetBuilding = &mockPlanetBuildingRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetBuildingRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
				},
			},
			"delete_planetResourceProductionRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetResourceProduction = &mockPlanetResourceProductionRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceProductionRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
				},
			},
			"delete_planetResourceStorageRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetResourceStorage = &mockPlanetResourceStorageRepository{
						errs: []error{errDefault},
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceStorageRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
				},
			},
			"delete_planetResourceRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.PlanetResource = &mockPlanetResourceRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetResourceRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteForPlanetCalled)
				},
			},
			"delete_planetRepositoryFails": {
				generateRepositoriesMocks: func() repositories.Repositories {
					repos := generatePlanetServiceMocks()
					repos.Planet = &mockPlanetRepository{
						err: errDefault,
					}

					return repos
				},
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
				expectedError: errDefault,
				verifyInteractions: func(repos repositories.Repositories, assert *require.Assertions) {
					m := assertPlanetRepoIsAMock(repos, assert)

					assert.Equal(1, m.deleteCalled)
				},
			},
		},

		returnTestCases: map[string]returnTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.Create(ctx, defaultPlanetDtoRequest)
					return out
				},
				expectedContent: communication.PlanetDtoResponse{
					Id:     defaultPlanet.Id,
					Player: defaultPlanet.Player,
					Name:   defaultPlanet.Name,

					CreatedAt: defaultPlanet.CreatedAt,
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.Get(ctx, defaultPlanetId)
					return out
				},
				expectedContent: communication.FullPlanetDtoResponse{
					PlanetDtoResponse: communication.PlanetDtoResponse{
						Id:     defaultPlanet.Id,
						Player: defaultPlanet.Player,
						Name:   defaultPlanet.Name,

						CreatedAt: defaultPlanet.CreatedAt,
					},
					Resources: []communication.PlanetResourceDtoResponse{
						{
							Planet:    defaultPlanet.Id,
							Resource:  metalResourceId,
							Amount:    784.0987,
							CreatedAt: testDate,
							UpdatedAt: testDate,
						},
					},
					Productions: []communication.PlanetResourceProductionDtoResponse{
						{
							Planet:     defaultPlanet.Id,
							Building:   &defaultBuilding.Id,
							Resource:   metalResourceId,
							Production: 31,
						},
					},
					Storages: []communication.PlanetResourceStorageDtoResponse{
						{
							Planet:   defaultPlanet.Id,
							Resource: metalResourceId,
							Storage:  9876,
						},
					},
					Buildings: []communication.PlanetBuildingDtoResponse{
						{
							Planet:    defaultPlanet.Id,
							Building:  defaultBuildingId,
							Level:     defaultPlanetBuilding.Level,
							CreatedAt: testDate,
							UpdatedAt: testDate,
						},
					},
					BuildingActions: []communication.BuildingActionDtoResponse{
						{
							Id:           defaultBuildingAction.Id,
							Planet:       defaultBuildingAction.Planet,
							Building:     defaultBuildingAction.Building,
							CurrentLevel: defaultBuildingAction.CurrentLevel,
							DesiredLevel: defaultBuildingAction.DesiredLevel,
							CreatedAt:    testDate,
							CompletedAt:  testDate,
						},
					},
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.List(ctx)
					return out
				},
				expectedContent: []communication.PlanetDtoResponse{
					{
						Id:     defaultPlanet.Id,
						Player: defaultPlanet.Player,
						Name:   defaultPlanet.Name,

						CreatedAt: defaultPlanet.CreatedAt,
					},
				},
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) interface{} {
					s := NewPlanetService(pool, repos)
					out, _ := s.ListForPlayer(ctx, defaultPlayerId)
					return out
				},
				expectedContent: []communication.PlanetDtoResponse{
					{
						Id:     defaultPlanet.Id,
						Player: defaultPlanet.Player,
						Name:   defaultPlanet.Name,

						CreatedAt: defaultPlanet.CreatedAt,
					},
				},
			},
		},

		transactionTestCases: map[string]transactionTestCase{
			"create": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Create(ctx, defaultPlanetDtoRequest)
					return err
				},
			},
			"get": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.Get(ctx, defaultPlanetId)
					return err
				},
			},
			"list": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.List(ctx)
					return err
				},
			},
			"listForPlayer": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					_, err := s.ListForPlayer(ctx, defaultPlayerId)
					return err
				},
			},
			"delete": {
				handler: func(ctx context.Context, pool db.ConnectionPool, repos repositories.Repositories) error {
					s := NewPlanetService(pool, repos)
					return s.Delete(ctx, defaultPlanetId)
				},
			},
		},
	}

	suite.Run(t, &s)
}

func generatePlanetServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		BuildingAction: &mockBuildingActionRepository{
			action: defaultBuildingAction,
		},
		BuildingActionCost: &mockBuildingActionCostRepository{
			actionCost: defaultBuildingActionCost,
		},
		BuildingActionResourceProduction: &mockBuildingActionResourceProductionRepository{
			actionResourceProduction: defaultBuildingActionResourceProduction,
		},
		Planet: &mockPlanetRepository{
			planet: defaultPlanet,
		},
		PlanetBuilding: &mockPlanetBuildingRepository{
			planetBuilding: defaultPlanetBuilding,
		},
		PlanetResource: &mockPlanetResourceRepository{
			planetResource: defaultPlanetResource,
		},
		PlanetResourceProduction: &mockPlanetResourceProductionRepository{
			planetResourceProduction: defaultPlanetResourceProduction,
		},
		PlanetResourceStorage: &mockPlanetResourceStorageRepository{
			planetResourceStorage: defaultPlanetResourceStorage,
		},
	}
}

func generateErrorPlanetServiceMocks() repositories.Repositories {
	return repositories.Repositories{
		Planet: &mockPlanetRepository{
			err: errDefault,
		},
	}
}

func assertPlanetRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetRepository {
	m, ok := repos.Planet.(*mockPlanetRepository)
	if !ok {
		assert.Fail("Provided planet repository is not a mock")
	}
	return m
}

func assertPlanetResourceRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetResourceRepository {
	m, ok := repos.PlanetResource.(*mockPlanetResourceRepository)
	if !ok {
		assert.Fail("Provided planet resource repository is not a mock")
	}
	return m
}

func assertPlanetResourceProductionRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetResourceProductionRepository {
	m, ok := repos.PlanetResourceProduction.(*mockPlanetResourceProductionRepository)
	if !ok {
		assert.Fail("Provided planet resource production repository is not a mock")
	}
	return m
}

func assertPlanetResourceStorageRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetResourceStorageRepository {
	m, ok := repos.PlanetResourceStorage.(*mockPlanetResourceStorageRepository)
	if !ok {
		assert.Fail("Provided planet resource storage repository is not a mock")
	}
	return m
}

func assertPlanetBuildingRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockPlanetBuildingRepository {
	m, ok := repos.PlanetBuilding.(*mockPlanetBuildingRepository)
	if !ok {
		assert.Fail("Provided planet building repository is not a mock")
	}
	return m
}

func assertBuildingActionRepoIsAMock(repos repositories.Repositories, assert *require.Assertions) *mockBuildingActionRepository {
	m, ok := repos.BuildingAction.(*mockBuildingActionRepository)
	if !ok {
		assert.Fail("Provided building action repository is not a mock")
	}
	return m
}
