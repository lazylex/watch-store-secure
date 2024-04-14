package joint

import (
	"context"
	"errors"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/in_memory"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/persistent"
	"log/slog"
)

type Repository struct {
	stateLocker StateLocker
	memory      in_memory.Interface
	persistent  persistent.Interface
}

func jointRepositoryError(text string) error {
	return errors.New("joint repo: " + text)
}

func New(memory in_memory.Interface, persistent persistent.Interface) Repository {
	go makeDataCache()
	return Repository{memory: memory, persistent: persistent, stateLocker: CreateStateLocker()}
}

// SaveSession сохраняет в памяти данные сессии
func (r *Repository) SaveSession(ctx context.Context, dto dto.SessionDTO) error {
	return r.memory.SaveSession(ctx, dto)
}

// SetAccountLoginData сохраняет в постоянном хранилище логин, хеш пароля, статус учетной записи и идентификатор
// пользователя. В памяти по возможности кеширует статус учетной записи, идентификатор пользователя и хеш пароля
func (r *Repository) SetAccountLoginData(ctx context.Context, data dto.AccountLoginDataDTO) error {
	var err error

	if err = r.persistent.SetAccountLoginData(ctx, data); err != nil {
		return err
	}

	if err = r.memory.SetAccountState(ctx, dto.LoginStateDTO{Login: data.Login, State: data.State}); err != nil {
		return err
	}

	r.memory.SetUserIdAndPasswordHash(ctx,
		dto.UserLoginAndIdWithPasswordHashDTO{UserId: data.UserId, Hash: data.Hash, Login: data.Login})
	return nil
}

// GetUserIdAndPasswordHash возвращает идентификатор пользователя и хеш его пароля
func (r *Repository) GetUserIdAndPasswordHash(ctx context.Context, login loginVO.Login) (dto.UserIdWithPasswordHashDTO, error) {
	idAndHash, err := r.memory.GetUserIdAndPasswordHash(ctx, login)
	if err == nil && idAndHash != (dto.UserIdWithPasswordHashDTO{}) {
		return idAndHash, err
	}
	data, errGetData := r.persistent.GetAccountLoginData(ctx, login)
	if errGetData != nil {
		return dto.UserIdWithPasswordHashDTO{}, errGetData
	}

	_ = r.saveToMemoryLoginData(ctx, data)

	return dto.UserIdWithPasswordHashDTO{UserId: data.UserId, Hash: data.Hash}, nil
}

// SetAccountState устанавливает состояние аккаунта
func (r *Repository) SetAccountState(ctx context.Context, stateDTO dto.LoginStateDTO) error {
	defer r.stateLocker.Unlock(stateDTO.Login)
	r.stateLocker.Lock(stateDTO.Login)

	if err := r.persistent.SetAccountState(ctx, stateDTO); err != nil {
		return err
	}

	return r.memory.SetAccountState(ctx, stateDTO)
}

// GetAccountState получает состояние учетной записи пользователя (сервиса)
func (r *Repository) GetAccountState(ctx context.Context, login loginVO.Login) (account_state.State, error) {
	var data dto.AccountLoginDataDTO
	var err error
	var state account_state.State

	r.stateLocker.WantRead(login)

	if state, err = r.memory.GetAccountStateByLogin(ctx, login); err == nil && account_state.IsStateCorrect(state) {
		return state, err
	}

	if data, err = r.GetAccountLoginData(ctx, login); err != nil {
		return 0, err
	}

	return data.State, nil
}

// GetAccountLoginData возвращает данные учетной записи по логину
func (r *Repository) GetAccountLoginData(ctx context.Context, login loginVO.Login) (dto.AccountLoginDataDTO, error) {
	var loginData dto.AccountLoginDataDTO
	var idAndHash dto.UserIdWithPasswordHashDTO
	var state account_state.State
	var err, errState, errHash error

	state, errState = r.memory.GetAccountStateByLogin(ctx, login)
	idAndHash, errHash = r.memory.GetUserIdAndPasswordHash(ctx, login)

	if errState == nil && errHash == nil && account_state.IsStateCorrect(state) {
		return dto.AccountLoginDataDTO{Login: login, UserId: idAndHash.UserId, Hash: idAndHash.Hash, State: state}, nil
	}

	if loginData, err = r.persistent.GetAccountLoginData(ctx, login); err != nil {
		return dto.AccountLoginDataDTO{}, err
	}

	_ = r.saveToMemoryLoginData(ctx, loginData)

	return loginData, nil
}

// AddService добавляет сервис в БД
func (r *Repository) AddService(ctx context.Context, data dto.NameWithDescriptionDTO) error {
	return r.persistent.AddService(ctx, data)
}

// AddPermission добавляет разрешение в БД
func (r *Repository) AddPermission(ctx context.Context, data dto.PermissionDTO) error {
	return r.persistent.AddPermission(ctx, data)
}

// AddRole добавляет роль в БД
func (r *Repository) AddRole(ctx context.Context, data dto.NameAndServiceWithDescriptionDTO) error {
	return r.persistent.AddRole(ctx, data)
}

// AddGroup добавляет группу в БД
func (r *Repository) AddGroup(ctx context.Context, data dto.NameAndServiceWithDescriptionDTO) error {
	return r.persistent.AddGroup(ctx, data)
}

// AssignRoleToGroup присоединяет роль к группе
func (r *Repository) AssignRoleToGroup(ctx context.Context, data dto.GroupRoleServiceNamesDTO) error {
	return r.persistent.AssignRoleToGroup(ctx, data)
}

// AssignRoleToAccount t назначает роль учетной записи
func (r *Repository) AssignRoleToAccount(ctx context.Context, data dto.RoleServiceNamesWithUserIdDTO) error {
	var err error
	if err = r.persistent.AssignRoleToAccount(ctx, data); err == nil && r.memory.ExistServicePermissionsNumbersForAccount(ctx, dto.ServiceNameWithUserIdDTO{
		UserId:  data.UserId,
		Service: data.Service,
	}) {
		r.refreshAccountPermissions(ctx, dto.ServiceNameWithUserIdDTO{UserId: data.UserId, Service: data.Service})
	}
	return err
}

// AssignGroupToAccount назначает группу учетной записи
func (r *Repository) AssignGroupToAccount(ctx context.Context, data dto.GroupServiceNamesWithUserIdDTO) error {
	var err error
	if err = r.persistent.AssignGroupToAccount(ctx, data); err == nil && r.memory.ExistServicePermissionsNumbersForAccount(ctx, dto.ServiceNameWithUserIdDTO{
		UserId:  data.UserId,
		Service: data.Service,
	}) {
		r.refreshAccountPermissions(ctx, dto.ServiceNameWithUserIdDTO{UserId: data.UserId, Service: data.Service})
	}
	return err
}

// AssignPermissionToRole назначает роли разрешение
func (r *Repository) AssignPermissionToRole(ctx context.Context, data dto.PermissionRoleServiceNamesDTO) error {
	return r.persistent.AssignPermissionToRole(ctx, data)
}

// AssignPermissionToGroup назначает разрешения группе
func (r *Repository) AssignPermissionToGroup(ctx context.Context, data dto.GroupPermissionServiceNamesDTO) error {
	return r.persistent.AssignPermissionToGroup(ctx, data)
}

// GetPermissionsForAccount возвращает название, номер и описание всех разрешений аккаунта для сервиса
func (r *Repository) GetPermissionsForAccount(ctx context.Context, data dto.ServiceNameWithUserIdDTO) ([]dto.PermissionWithoutServiceDTO, error) {
	return r.persistent.GetPermissionsForAccount(ctx, data)
}

// GetPermissionsNumbersForAccount возвращает номера всех разрешений аккаунта для сервиса
func (r *Repository) GetPermissionsNumbersForAccount(ctx context.Context, data dto.ServiceNameWithUserIdDTO) ([]int, error) {
	return r.persistent.GetPermissionsNumbersForAccount(ctx, data)
}

// saveToMemoryLoginData сохраняет в памяти данные, необходимые для процесса входа в систему пользователя (сервиса)
func (r *Repository) saveToMemoryLoginData(ctx context.Context, data dto.AccountLoginDataDTO) error {
	r.memory.SetUserIdAndPasswordHash(ctx, dto.UserLoginAndIdWithPasswordHashDTO{UserId: data.UserId, Login: data.Login, Hash: data.Hash})
	return r.memory.SetAccountState(ctx, dto.LoginStateDTO{Login: data.Login, State: data.State})
}

// refreshAccountPermissions обновляет кеш разрешений
func (r *Repository) refreshAccountPermissions(ctx context.Context, data dto.ServiceNameWithUserIdDTO) {
	// TODO implement
	slog.Debug("not implemented permission refresh")
}

// makeDataCache считывает все данные (которые возможно кешировать) из постоянного хранилища в хранилище в памяти
func makeDataCache() {
	// TODO implement
	slog.Debug(jointRepositoryError("makeDataCache not implemented").Error())
}
