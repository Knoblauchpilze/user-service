package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/galactic-sovereign/pkg/db"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/persistence"
	"github.com/KnoblauchPilze/galactic-sovereign/pkg/repositories"
	"github.com/google/uuid"
)

type mockAclRepository struct {
	repositories.AclRepository

	aclIds []uuid.UUID
	acl    persistence.Acl

	getErr        error
	getForUserErr error
	deleteErr     error

	inAclIds         []uuid.UUID
	getCalled        int
	inUserId         uuid.UUID
	getForUserCalled int
	deleteCalled     int
}

func (m *mockAclRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Acl, error) {
	m.getCalled++
	m.inAclIds = append(m.inAclIds, id)
	return m.acl, m.getErr
}

func (m *mockAclRepository) GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserCalled++
	m.inUserId = user
	return m.aclIds, m.getForUserErr
}

func (m *mockAclRepository) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	m.deleteCalled++
	m.inUserId = user
	return m.deleteErr
}

type mockApiKeyRepository struct {
	repositories.ApiKeyRepository

	apiKey    persistence.ApiKey
	apiKeyIds []uuid.UUID
	createErr error
	getErr    error
	deleteErr error

	createCalled        int
	createdApiKey       persistence.ApiKey
	getForKeyCalled     int
	apiKeyId            uuid.UUID
	getForUserCalled    int
	userId              uuid.UUID
	deleteCalled        int
	deleteIds           []uuid.UUID
	deleteForUserCalled int
	deleteUserId        uuid.UUID
}

func (m *mockApiKeyRepository) Create(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error) {
	m.createCalled++
	m.createdApiKey = apiKey
	return apiKey, m.createErr
}

func (m *mockApiKeyRepository) GetForUser(ctx context.Context, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserCalled++
	m.userId = user
	return m.apiKeyIds, m.getErr
}

func (m *mockApiKeyRepository) GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error) {
	m.getForKeyCalled++
	m.apiKeyId = apiKey
	return m.apiKey, m.getErr
}

func (m *mockApiKeyRepository) Delete(ctx context.Context, ids []uuid.UUID) error {
	m.deleteCalled++
	m.deleteIds = ids
	return m.deleteErr
}

func (m *mockApiKeyRepository) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	m.deleteForUserCalled++
	m.deleteUserId = user
	return m.deleteErr
}

type mockBuildingRepository struct {
	repositories.BuildingRepository

	building persistence.Building
	err      error

	listCalled int
}

func (m *mockBuildingRepository) List(ctx context.Context, tx db.Transaction) ([]persistence.Building, error) {
	m.listCalled++
	return []persistence.Building{m.building}, m.err
}

type mockBuildingActionRepository struct {
	repositories.BuildingActionRepository

	action persistence.BuildingAction
	errs   []error
	calls  int

	createCalled                   int
	createdBuildingAction          persistence.BuildingAction
	getCalled                      int
	getId                          uuid.UUID
	listForPlanetId                uuid.UUID
	listForPlanetCalled            int
	listBeforeCompletionTimeCalled int
	listBeforeCompletionTimePlanet uuid.UUID
	listBeforeCompletionTime       time.Time
	deleteCalled                   int
	deleteId                       uuid.UUID
	deleteForPlanetCalled          int
	deleteForPlanetId              uuid.UUID
}

func (m *mockBuildingActionRepository) Create(ctx context.Context, tx db.Transaction, action persistence.BuildingAction) (persistence.BuildingAction, error) {
	m.createCalled++
	m.createdBuildingAction = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.action, *err
}

func (m *mockBuildingActionRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.BuildingAction, error) {
	m.getCalled++
	m.getId = id

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.action, *err
}

func (m *mockBuildingActionRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.BuildingAction, error) {
	m.listForPlanetCalled++
	m.listForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.BuildingAction{m.action}, *err
}

func (m *mockBuildingActionRepository) ListBeforeCompletionTime(ctx context.Context, tx db.Transaction, planet uuid.UUID, until time.Time) ([]persistence.BuildingAction, error) {
	m.listBeforeCompletionTimeCalled++
	m.listBeforeCompletionTimePlanet = planet
	m.listBeforeCompletionTime = until

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.BuildingAction{m.action}, *err
}

func (m *mockBuildingActionRepository) Delete(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

func (m *mockBuildingActionRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

type mockBuildingActionCostRepository struct {
	repositories.BuildingActionCostRepository

	actionCost persistence.BuildingActionCost
	errs       []error
	calls      int

	createCalled              int
	createdBuildingActionCost persistence.BuildingActionCost
	listForActionCalled       int
	listForActionId           uuid.UUID
	deleteForActionCalled     int
	deleteForActionId         uuid.UUID
	deleteForPlanetCalled     int
	deleteForPlanetId         uuid.UUID
}

func (m *mockBuildingActionCostRepository) Create(ctx context.Context, tx db.Transaction, cost persistence.BuildingActionCost) (persistence.BuildingActionCost, error) {
	m.createCalled++
	m.createdBuildingActionCost = cost

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.actionCost, *err
}

func (m *mockBuildingActionCostRepository) ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionCost, error) {
	m.listForActionCalled++
	m.listForActionId = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.BuildingActionCost{m.actionCost}, *err
}

func (m *mockBuildingActionCostRepository) DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	m.deleteForActionCalled++
	m.deleteForActionId = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

func (m *mockBuildingActionCostRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

type mockBuildingActionResourceProductionRepository struct {
	repositories.BuildingActionResourceProductionRepository

	actionResourceProduction persistence.BuildingActionResourceProduction
	errs                     []error
	calls                    int

	createCalled                            int
	createdBuildingActionResourceProduction persistence.BuildingActionResourceProduction
	listForActionCalled                     int
	listForActionId                         uuid.UUID
	deleteForActionCalled                   int
	deleteForActionId                       uuid.UUID
	deleteForPlanetCalled                   int
	deleteForPlanetId                       uuid.UUID
}

func (m *mockBuildingActionResourceProductionRepository) Create(ctx context.Context, tx db.Transaction, production persistence.BuildingActionResourceProduction) (persistence.BuildingActionResourceProduction, error) {
	m.createCalled++
	m.createdBuildingActionResourceProduction = production

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.actionResourceProduction, *err
}

func (m *mockBuildingActionResourceProductionRepository) ListForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) ([]persistence.BuildingActionResourceProduction, error) {
	m.listForActionCalled++
	m.listForActionId = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.BuildingActionResourceProduction{m.actionResourceProduction}, *err
}

func (m *mockBuildingActionResourceProductionRepository) DeleteForAction(ctx context.Context, tx db.Transaction, action uuid.UUID) error {
	m.deleteForActionCalled++
	m.deleteForActionId = action

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

func (m *mockBuildingActionResourceProductionRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

type mockBuildingCostRepository struct {
	repositories.BuildingCostRepository

	buildingCost persistence.BuildingCost
	err          error

	listForBuildingId     uuid.UUID
	listForBuildingCalled int
}

func (m *mockBuildingCostRepository) ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingCost, error) {
	m.listForBuildingCalled++
	m.listForBuildingId = building
	return []persistence.BuildingCost{m.buildingCost}, m.err
}

type mockBuildingResourceProductionRepository struct {
	repositories.BuildingResourceProductionRepository

	buildingResourceProduction persistence.BuildingResourceProduction
	err                        error

	listForBuildingId     uuid.UUID
	listForBuildingCalled int
}

func (m *mockBuildingResourceProductionRepository) ListForBuilding(ctx context.Context, tx db.Transaction, building uuid.UUID) ([]persistence.BuildingResourceProduction, error) {
	m.listForBuildingCalled++
	m.listForBuildingId = building
	return []persistence.BuildingResourceProduction{m.buildingResourceProduction}, m.err
}

type mockPlanetBuildingRepository struct {
	repositories.PlanetBuildingRepository

	planetBuilding persistence.PlanetBuilding
	err            error
	updateErr      error

	getForPlanetAndBuildingCalled   int
	getForPlanetAndBuildingPlanet   uuid.UUID
	getForPlanetAndBuildingBuilding uuid.UUID
	listForPlanetCalled             int
	listForPlanetId                 uuid.UUID
	updateCalled                    int
	updateBuilding                  persistence.PlanetBuilding
	deleteForPlanetCalled           int
	deleteForPlanetId               uuid.UUID
}

func (m *mockPlanetBuildingRepository) GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building uuid.UUID) (persistence.PlanetBuilding, error) {
	m.getForPlanetAndBuildingCalled++
	m.getForPlanetAndBuildingPlanet = planet
	m.getForPlanetAndBuildingBuilding = building
	return m.planetBuilding, m.err
}

func (m *mockPlanetBuildingRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetBuilding, error) {
	m.listForPlanetCalled++
	m.listForPlanetId = planet
	return []persistence.PlanetBuilding{m.planetBuilding}, m.err
}

func (m *mockPlanetBuildingRepository) Update(ctx context.Context, tx db.Transaction, building persistence.PlanetBuilding) (persistence.PlanetBuilding, error) {
	m.updateCalled++
	m.updateBuilding = building
	return m.updateBuilding, m.updateErr
}

func (m *mockPlanetBuildingRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet
	return m.err
}

type mockPlanetRepository struct {
	repositories.PlanetRepository

	planet persistence.Planet
	err    error

	createCalled        int
	createdPlanet       persistence.Planet
	getCalled           int
	getId               uuid.UUID
	listCalled          int
	listForPlayerId     uuid.UUID
	listForPlayerCalled int
	deleteCalled        int
	deleteId            uuid.UUID
}

func (m *mockPlanetRepository) Create(ctx context.Context, tx db.Transaction, planet persistence.Planet) (persistence.Planet, error) {
	m.createCalled++
	m.createdPlanet = planet
	return m.planet, m.err
}

func (m *mockPlanetRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Planet, error) {
	m.getCalled++
	m.getId = id
	return m.planet, m.err
}

func (m *mockPlanetRepository) List(ctx context.Context, tx db.Transaction) ([]persistence.Planet, error) {
	m.listCalled++
	return []persistence.Planet{m.planet}, m.err
}

func (m *mockPlanetRepository) ListForPlayer(ctx context.Context, tx db.Transaction, player uuid.UUID) ([]persistence.Planet, error) {
	m.listForPlayerCalled++
	m.listForPlayerId = player
	return []persistence.Planet{m.planet}, m.err
}

func (m *mockPlanetRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

type mockPlanetResourceRepository struct {
	repositories.PlanetResourceRepository

	planetResource persistence.PlanetResource
	err            error
	updateErr      error

	createCalled           int
	createdPlanetResource  persistence.PlanetResource
	listForPlanetIds       []uuid.UUID
	listForPlanetCalled    int
	updateCalled           int
	updatedPlanetResources []persistence.PlanetResource
	deleteForPlanetCalled  int
	deleteForPlanetId      uuid.UUID
}

func (m *mockPlanetResourceRepository) Create(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	m.createCalled++
	m.createdPlanetResource = resource
	return m.planetResource, m.err
}

func (m *mockPlanetResourceRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResource, error) {
	m.listForPlanetCalled++
	m.listForPlanetIds = append(m.listForPlanetIds, planet)
	return []persistence.PlanetResource{m.planetResource}, m.err
}

func (m *mockPlanetResourceRepository) Update(ctx context.Context, tx db.Transaction, resource persistence.PlanetResource) (persistence.PlanetResource, error) {
	m.updateCalled++
	m.updatedPlanetResources = append(m.updatedPlanetResources, resource)
	return resource, m.updateErr
}

func (m *mockPlanetResourceRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet
	return m.err
}

type mockPlanetResourceProductionRepository struct {
	repositories.PlanetResourceProductionRepository

	planetResourceProduction persistence.PlanetResourceProduction
	errs                     []error
	calls                    int
	updateErr                error

	createCalled                     int
	createdPlanetResourceProduction  persistence.PlanetResourceProduction
	getForPlanetAndBuildingCalled    int
	getForPlanetAndBuildingPlanet    uuid.UUID
	getForPlanetAndBuildingBuilding  *uuid.UUID
	listForPlanetIds                 []uuid.UUID
	listForPlanetCalled              int
	updateCalled                     int
	updatedPlanetResourceProductions []persistence.PlanetResourceProduction
	deleteForPlanetCalled            int
	deleteForPlanetId                uuid.UUID
}

func (m *mockPlanetResourceProductionRepository) Create(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	m.createCalled++
	m.createdPlanetResourceProduction = production

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.planetResourceProduction, *err
}

func (m *mockPlanetResourceProductionRepository) GetForPlanetAndBuilding(ctx context.Context, tx db.Transaction, planet uuid.UUID, building *uuid.UUID) (persistence.PlanetResourceProduction, error) {
	m.getForPlanetAndBuildingCalled++
	m.getForPlanetAndBuildingPlanet = planet
	m.getForPlanetAndBuildingBuilding = building

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.planetResourceProduction, *err
}

func (m *mockPlanetResourceProductionRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceProduction, error) {
	m.listForPlanetCalled++
	m.listForPlanetIds = append(m.listForPlanetIds, planet)

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.PlanetResourceProduction{m.planetResourceProduction}, *err
}

func (m *mockPlanetResourceProductionRepository) Update(ctx context.Context, tx db.Transaction, production persistence.PlanetResourceProduction) (persistence.PlanetResourceProduction, error) {
	m.updateCalled++
	m.updatedPlanetResourceProductions = append(m.updatedPlanetResourceProductions, production)
	return production, m.updateErr
}

func (m *mockPlanetResourceProductionRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

type mockPlanetResourceStorageRepository struct {
	repositories.PlanetResourceStorageRepository

	planetResourceStorage persistence.PlanetResourceStorage
	errs                  []error
	calls                 int
	updateErr             error

	createCalled                  int
	createdPlanetResourceStorage  persistence.PlanetResourceStorage
	listForPlanetIds              []uuid.UUID
	listForPlanetCalled           int
	updateCalled                  int
	updatedPlanetResourceStorages []persistence.PlanetResourceStorage
	deleteForPlanetCalled         int
	deleteForPlanetId             uuid.UUID
}

func (m *mockPlanetResourceStorageRepository) Create(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	m.createCalled++
	m.createdPlanetResourceStorage = storage

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return m.planetResourceStorage, *err
}

func (m *mockPlanetResourceStorageRepository) ListForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) ([]persistence.PlanetResourceStorage, error) {
	m.listForPlanetCalled++
	m.listForPlanetIds = append(m.listForPlanetIds, planet)

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return []persistence.PlanetResourceStorage{m.planetResourceStorage}, *err
}

func (m *mockPlanetResourceStorageRepository) Update(ctx context.Context, tx db.Transaction, storage persistence.PlanetResourceStorage) (persistence.PlanetResourceStorage, error) {
	m.updateCalled++
	m.updatedPlanetResourceStorages = append(m.updatedPlanetResourceStorages, storage)
	return storage, m.updateErr
}

func (m *mockPlanetResourceStorageRepository) DeleteForPlanet(ctx context.Context, tx db.Transaction, planet uuid.UUID) error {
	m.deleteForPlanetCalled++
	m.deleteForPlanetId = planet

	err := getValueToReturnOr(m.calls, m.errs, nil)
	m.calls++

	return *err
}

type mockPlayerRepository struct {
	repositories.PlayerRepository

	player persistence.Player
	err    error

	createCalled         int
	createdPlayer        persistence.Player
	getCalled            int
	getId                uuid.UUID
	listCalled           int
	listForApiUserId     uuid.UUID
	listForApiUserCalled int
	deleteCalled         int
	deleteId             uuid.UUID
}

func (m *mockPlayerRepository) Create(ctx context.Context, tx db.Transaction, player persistence.Player) (persistence.Player, error) {
	m.createCalled++
	m.createdPlayer = player
	return m.player, m.err
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uuid.UUID) (persistence.Player, error) {
	m.getCalled++
	m.getId = id
	return m.player, m.err
}

func (m *mockPlayerRepository) List(ctx context.Context) ([]persistence.Player, error) {
	m.listCalled++
	return []persistence.Player{m.player}, m.err
}

func (m *mockPlayerRepository) ListForApiUser(ctx context.Context, apiUser uuid.UUID) ([]persistence.Player, error) {
	m.listForApiUserCalled++
	m.listForApiUserId = apiUser
	return []persistence.Player{m.player}, m.err
}

func (m *mockPlayerRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

type mockResourceRepository struct {
	repositories.ResourceRepository

	resources []persistence.Resource
	err       error

	listCalled int
}

func (m *mockResourceRepository) List(ctx context.Context, tx db.Transaction) ([]persistence.Resource, error) {
	m.listCalled++
	return m.resources, m.err
}

type mockUniverseRepository struct {
	repositories.UniverseRepository

	universe persistence.Universe
	err      error

	createCalled    int
	createdUniverse persistence.Universe
	getCalled       int
	getId           uuid.UUID
	listCalled      int
	deleteCalled    int
	deleteId        uuid.UUID
}

func (m *mockUniverseRepository) Create(ctx context.Context, universe persistence.Universe) (persistence.Universe, error) {
	m.createCalled++
	m.createdUniverse = universe
	return m.universe, m.err
}

func (m *mockUniverseRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Universe, error) {
	m.getCalled++
	m.getId = id
	return m.universe, m.err
}

func (m *mockUniverseRepository) List(ctx context.Context) ([]persistence.Universe, error) {
	m.listCalled++
	return []persistence.Universe{m.universe}, m.err
}

func (m *mockUniverseRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

type mockUserRepository struct {
	repositories.UserRepository

	user      persistence.User
	ids       []uuid.UUID
	err       error
	updateErr error

	createCalled   int
	createdUser    persistence.User
	getCalled      int
	getId          uuid.UUID
	getEmailCalled int
	getEmail       string
	listCalled     int
	updateCalled   int
	updatedUser    persistence.User
	deleteCalled   int
	deleteId       uuid.UUID
}

func (m *mockUserRepository) Create(ctx context.Context, user persistence.User) (persistence.User, error) {
	m.createCalled++
	m.createdUser = user
	return m.user, m.err
}

func (m *mockUserRepository) Get(ctx context.Context, id uuid.UUID) (persistence.User, error) {
	m.getCalled++
	m.getId = id
	return m.user, m.err
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (persistence.User, error) {
	m.getEmailCalled++
	m.getEmail = email
	return m.user, m.err
}

func (m *mockUserRepository) List(ctx context.Context) ([]uuid.UUID, error) {
	m.listCalled++
	return m.ids, m.err
}

func (m *mockUserRepository) Update(ctx context.Context, user persistence.User) (persistence.User, error) {
	m.updateCalled++
	m.updatedUser = user
	return m.updatedUser, m.updateErr
}

func (m *mockUserRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

type mockUserLimitRepository struct {
	repositories.UserLimitRepository

	userLimitIds []uuid.UUID
	userLimit    persistence.UserLimit

	getErr        error
	getForUserErr error
	deleteErr     error

	inUserLimitIds   []uuid.UUID
	getCalled        int
	inUserId         uuid.UUID
	getForUserCalled int
	deleteCalled     int
}

func (m *mockUserLimitRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.UserLimit, error) {
	m.getCalled++
	m.inUserLimitIds = append(m.inUserLimitIds, id)
	return m.userLimit, m.getErr
}

func (m *mockUserLimitRepository) GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserCalled++
	m.inUserId = user
	return m.userLimitIds, m.getForUserErr
}

func (m *mockUserLimitRepository) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	m.deleteCalled++
	m.inUserId = user
	return m.deleteErr
}

func getValueToReturnOr[T any](count int, values []T, value T) *T {
	out := getValueToReturn(count, values)
	if out == nil {
		return &value
	}

	return out
}

func getValueToReturn[T any](count int, values []T) *T {
	var out *T
	if count > len(values) {
		count = 0
	}
	if count < len(values) {
		out = &values[count]
	}

	return out
}
