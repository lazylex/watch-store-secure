/*
Package joint: пакет для работы с объединенным хранилищем, служащим в качестве единой точки обращения к in memory и
постоянному хранилищу. Для создания структуры Repository, с методами для работы объединенного хранилища, используется
функция MustCreate. Эта же функция в фоне проводит кеширование части данных из постоянного хранилища в оперативную
память.
*/
package joint

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/errors/joint"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/in_memory"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/persistent"
	"log/slog"
	"os"
	"time"
)

// Repository структура, содержащая блокировщик состояния и объекты, реализующие интерфейсы работы с постоянным и
// in memory хранилищами.
type Repository struct {
	stateLocker StateLocker          // Структура, управляющая блокировками чтения/записи состояния аккаунтов
	memory      in_memory.Interface  // Объект, реализующий работу с хранилищем, расположенном в ОЗУ
	persistent  persistent.Interface // Объект, реализующий работу с постоянным хранилищем
}

// MustCreate создает объединенное хранилище из in memory кеша и постоянного хранилища. Возвращает структуру Repository
// для взаимодействия с хранилищем, параллельно проводя кеширование части данных из постоянного хранилища в память.
func MustCreate(memoryRepo in_memory.Interface, persistentRepo persistent.Interface) Repository {
	if memoryRepo == nil || persistentRepo == nil {
		slog.Error(adaptErr(joint.ErrNilRepo).Error())
		os.Exit(1)
	}
	r := Repository{memory: memoryRepo, persistent: persistentRepo, stateLocker: CreateStateLocker()}
	go r.makeDataCache()
	return r
}

// SaveSession сохраняет в памяти данные сессии.
func (r *Repository) SaveSession(ctx context.Context, dto *dto.UserIdToken) error {
	return adaptErr(r.memory.SaveSession(ctx, dto))
}

// DeleteSession удаляет сессию.
func (r *Repository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	return adaptErr(r.memory.DeleteSession(ctx, id))
}

// UserUUIDFromSession получает UUID пользователя сессии.
func (r *Repository) UserUUIDFromSession(ctx context.Context, sessionToken string) (uuid.UUID, error) {
	id, err := r.memory.UserUUIDFromSession(ctx, sessionToken)
	return id, adaptErr(err)
}

// SetAccountLoginData сохраняет в постоянном хранилище логин, хеш пароля, статус учетной записи и идентификатор
// пользователя. В памяти по возможности кеширует статус учетной записи, идентификатор пользователя и хеш пароля.
func (r *Repository) SetAccountLoginData(ctx context.Context, data *dto.UserIdLoginHashState) error {
	var err error

	if err = r.persistent.SetAccountLoginData(ctx, data); err != nil {
		return adaptErr(err)
	}

	ctx = context.Background()

	if err = r.memory.SetAccountState(ctx, &dto.LoginState{Login: data.Login, State: data.State}); err != nil {
		return adaptErr(err)
	}

	r.memory.SetUserIdAndPasswordHash(ctx,
		&dto.UserIdLoginHash{UserId: data.UserId, Hash: data.Hash, Login: data.Login})
	return nil
}

// UserIdAndPasswordHash возвращает идентификатор пользователя и хеш его пароля.
func (r *Repository) UserIdAndPasswordHash(ctx context.Context, login loginVO.Login) (dto.UserIdHash, error) {
	idAndHash, err := r.memory.UserIdAndPasswordHash(ctx, login)
	if err == nil && idAndHash != (dto.UserIdHash{}) {
		return idAndHash, err
	}
	data, errGetData := r.persistent.AccountLoginData(ctx, login)
	if errGetData != nil {
		return dto.UserIdHash{}, adaptErr(errGetData)
	}

	_ = r.saveToMemoryLoginData(context.Background(), &data)

	return dto.UserIdHash{UserId: data.UserId, Hash: data.Hash}, nil
}

// SessionToken возвращает токен сессии по идентификатору пользователя.
func (r *Repository) SessionToken(ctx context.Context, id uuid.UUID) (string, error) {
	token, err := r.memory.SessionToken(ctx, id)
	return token, adaptErr(err)
}

// SetAccountState устанавливает состояние аккаунта.
func (r *Repository) SetAccountState(ctx context.Context, data *dto.LoginState) error {
	defer r.stateLocker.Unlock(data.Login)
	r.stateLocker.Lock(data.Login)

	if err := r.persistent.SetAccountState(ctx, data); err != nil {
		return adaptErr(err)
	}

	return adaptErr(r.memory.SetAccountState(ctx, data))
}

// AccountState получает состояние учетной записи пользователя (сервиса).
func (r *Repository) AccountState(ctx context.Context, login loginVO.Login) (account_state.State, error) {
	var data dto.UserIdLoginHashState
	var err error
	var state account_state.State

	r.stateLocker.WantRead(login)

	if state, err = r.memory.AccountStateByLogin(ctx, login); err == nil {
		return state, adaptErr(err)
	}

	if data, err = r.AccountLoginData(ctx, login); err != nil {
		return 0, adaptErr(err)
	}

	return data.State, nil
}

// AccountLoginData возвращает данные учетной записи по логину.
func (r *Repository) AccountLoginData(ctx context.Context, login loginVO.Login) (dto.UserIdLoginHashState, error) {
	var loginData dto.UserIdLoginHashState
	var idAndHash dto.UserIdHash
	var state account_state.State
	var err, errState, errHash error

	state, errState = r.memory.AccountStateByLogin(ctx, login)
	idAndHash, errHash = r.memory.UserIdAndPasswordHash(ctx, login)

	if errState == nil && errHash == nil {
		return dto.UserIdLoginHashState{Login: login, UserId: idAndHash.UserId, Hash: idAndHash.Hash, State: state}, nil
	}

	if loginData, err = r.persistent.AccountLoginData(ctx, login); err != nil {
		return dto.UserIdLoginHashState{}, adaptErr(err)
	}

	_ = r.saveToMemoryLoginData(context.Background(), &loginData)

	return loginData, nil
}

// CreateService добавляет сервис в БД.
func (r *Repository) CreateService(ctx context.Context, data *dto.NameDescription) error {
	return adaptErr(r.persistent.CreateService(ctx, data))
}

// CreatePermission добавляет разрешение в БД.
func (r *Repository) CreatePermission(ctx context.Context, data *dto.NameServiceDescription) error {
	return adaptErr(r.persistent.CreatePermission(ctx, data))
}

// CreateRole добавляет роль в БД.
func (r *Repository) CreateRole(ctx context.Context, data *dto.NameServiceDescription) error {
	return adaptErr(r.persistent.CreateRole(ctx, data))
}

// CreateGroup добавляет группу в БД.
func (r *Repository) CreateGroup(ctx context.Context, data *dto.NameServiceDescription) error {
	return adaptErr(r.persistent.CreateGroup(ctx, data))
}

// CreateOrUpdateInstance добавляет/обновляет в БД название экземпляра сервиса и секретный ключ для подписи токена.
func (r *Repository) CreateOrUpdateInstance(ctx context.Context, data *dto.NameServiceSecret) error {
	if err := r.persistent.CreateOrUpdateInstance(ctx, data); err != nil {
		return adaptErr(err)
	}

	if err := r.memory.SetInstanceServiceAndSecret(context.Background(), data); err != nil {
		return adaptErr(joint.ErrCacheSavedData)
	}

	return nil
}

// AssignRoleToGroup присоединяет роль к группе.
func (r *Repository) AssignRoleToGroup(ctx context.Context, data *dto.GroupRoleService) error {
	return adaptErr(r.persistent.AssignRoleToGroup(ctx, data))
}

// AssignRoleToAccount назначает роль учетной записи.
func (r *Repository) AssignRoleToAccount(ctx context.Context, data *dto.UserIdRoleService) error {
	var err error
	if err = r.persistent.AssignRoleToAccount(ctx, data); err == nil && r.memory.ExistServicePermissionsNumbersForAccount(ctx, &dto.UserIdService{
		UserId:  data.UserId,
		Service: data.Service,
	}) {
		r.refreshAccountPermissions(context.Background(), &dto.UserIdService{UserId: data.UserId, Service: data.Service})
	}
	return adaptErr(err)
}

// AssignGroupToAccount назначает группу учетной записи.
func (r *Repository) AssignGroupToAccount(ctx context.Context, data *dto.UserIdGroupService) error {
	var err error
	if err = r.persistent.AssignGroupToAccount(ctx, data); err == nil && r.memory.ExistServicePermissionsNumbersForAccount(ctx, &dto.UserIdService{
		UserId:  data.UserId,
		Service: data.Service,
	}) {
		r.refreshAccountPermissions(context.Background(), &dto.UserIdService{UserId: data.UserId, Service: data.Service})
	}
	return adaptErr(err)
}

// AssignInstancePermissionToAccount прикрепляет разрешение конкретного экземпляра сервиса к учетной записи.
func (r *Repository) AssignInstancePermissionToAccount(ctx context.Context, data *dto.UserIdInstancePermission) error {
	var err error
	var number int

	if err = r.persistent.AssignInstancePermissionToAccount(ctx, data); err != nil {
		return adaptErr(err)
	}

	if number, err = r.persistent.PermissionNumber(ctx, data.Permission, data.Instance); err != nil {
		return adaptErr(joint.ErrCacheSavedData)
	}

	if err = r.memory.SetInstancePermissionsNumbersForAccount(context.Background(), &dto.UserIdInstancePermNumbers{
		UserId:            data.UserId,
		Instance:          data.Instance,
		PermissionNumbers: []int{number},
	}); err != nil {
		return adaptErr(joint.ErrCacheSavedData)
	}

	return nil
}

// AssignPermissionToRole назначает роли разрешение.
func (r *Repository) AssignPermissionToRole(ctx context.Context, data *dto.PermissionRoleService) error {
	return adaptErr(r.persistent.AssignPermissionToRole(ctx, data))
}

// AssignPermissionToGroup назначает разрешения группе.
func (r *Repository) AssignPermissionToGroup(ctx context.Context, data *dto.GroupPermissionService) error {
	return adaptErr(r.persistent.AssignPermissionToGroup(ctx, data))
}

// ServicePermissionsForAccount возвращает название, номер и описание разрешений аккаунта для сервиса.
func (r *Repository) ServicePermissionsForAccount(ctx context.Context, data *dto.UserIdService) ([]dto.NameNumberDescription, error) {
	permissions, err := r.persistent.ServicePermissionsForAccount(ctx, data)
	return permissions, adaptErr(err)
}

// ServicePermissionsNumbersForAccount возвращает номера разрешений аккаунта для сервиса.
func (r *Repository) ServicePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdService) ([]int, error) {
	var numbers []int
	var err error

	if numbers, err = r.memory.ServicePermissionsNumbersForAccount(ctx, data); err != nil || len(numbers) == 0 {
		numbers, err = r.servicePermissionsNumbersForAccountFromPersistentWithSaveToMemory(ctx, data)
	}

	return numbers, adaptErr(err)
}

// ServiceNumberedPermissions возвращает номера и названия разрешений сервиса.
func (r *Repository) ServiceNumberedPermissions(ctx context.Context, serviceName string) (*[]dto.NameNumber, error) {
	var err error
	var result *[]dto.NameNumber

	if result, err = r.memory.ServiceNumberedPermissions(ctx, serviceName); err == nil && result != nil && len(*result) != 0 {
		return result, nil
	}

	if result, err = r.persistent.ServiceNumberedPermissions(ctx, serviceName); err != nil {
		return nil, adaptErr(err)
	}

	go func() {
		_ = r.memory.SetServiceNumberedPermissions(context.Background(), serviceName, result)
	}()

	return result, nil
}

// servicePermissionsNumbersForAccountFromPersistentWithSaveToMemory возвращает номера разрешений аккаунта для сервиса и
// кеширует их в память.
func (r *Repository) servicePermissionsNumbersForAccountFromPersistentWithSaveToMemory(ctx context.Context, data *dto.UserIdService) ([]int, error) {
	var numbers []int
	var err error

	if numbers, err = r.persistent.ServicePermissionsNumbersForAccount(ctx, data); err == nil && len(numbers) > 0 {
		go func() {
			_ = r.memory.SetServicePermissionsNumbersForAccount(context.Background(), &dto.UserIdServicePermNumbers{
				UserId:            data.UserId,
				Service:           data.Service,
				PermissionNumbers: numbers,
			})
		}()
	}

	return numbers, adaptErr(err)
}

// InstancePermissionsNumbersForAccount возвращает номера разрешений аккаунта для экземпляра сервиса.
func (r *Repository) InstancePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdInstance) ([]int, error) {
	var numbers []int
	var err error

	if numbers, err = r.memory.InstancePermissionsNumbersForAccount(ctx, data); err != nil || len(numbers) == 0 {
		numbers, err = r.instancePermissionsNumbersForAccountFromPersistentWithSaveToMemory(ctx, data)
	}

	return numbers, adaptErr(err)
}

// ServiceName возвращает название сервиса переданного экземпляра.
func (r *Repository) ServiceName(ctx context.Context, instanceName string) (string, error) {
	var name string
	var err error

	if name, err = r.memory.ServiceName(ctx, instanceName); err == nil {
		return name, nil
	}

	if name, err = r.persistent.ServiceName(ctx, instanceName); err != nil {
		return "", adaptErr(err)
	}

	defer func() {
		_ = r.memory.SetInstanceServiceName(context.Background(), &dto.NameService{Name: instanceName, Service: name})
	}()

	return name, nil
}

// ServicesNames возвращает список названий всех сервисов.
func (r *Repository) ServicesNames(ctx context.Context) ([]string, error) {
	result, err := r.persistent.ServicesNames(ctx)
	return result, adaptErr(err)
}

// InstanceSecret возвращает строку, необходимую для подписи токена, предназначенного для взаимодействия с
// соответствующим экземпляром сервиса.
func (r *Repository) InstanceSecret(ctx context.Context, name string) (string, error) {
	var secret string
	var err error
	if secret, err = r.memory.InstanceSecret(ctx, name); err == nil {
		return secret, nil
	}

	if secret, err = r.persistent.InstanceSecret(ctx, name); err != nil {
		return "", adaptErr(err)
	}

	defer func() {
		_ = r.memory.SetInstanceSecret(context.Background(), &dto.NameSecret{Name: name, Secret: secret})
	}()

	return secret, err
}

// DeleteRole удаляет роль из БД.
func (r *Repository) DeleteRole(ctx context.Context, data *dto.NameService) error {
	return adaptErr(r.persistent.DeleteRole(ctx, data))
}

// DeleteGroup удаляет группу из БД.
func (r *Repository) DeleteGroup(ctx context.Context, data *dto.NameService) error {
	return adaptErr(r.persistent.DeleteGroup(ctx, data))
}

// DeletePermission удаляет разрешение из БД.
func (r *Repository) DeletePermission(ctx context.Context, data *dto.NameService) error {
	return adaptErr(r.persistent.DeletePermission(ctx, data))
}

// instancePermissionsNumbersForAccountFromPersistentWithSaveToMemory возвращает номера разрешений аккаунта для
// экземпляра сервиса и кеширует их в память.
func (r *Repository) instancePermissionsNumbersForAccountFromPersistentWithSaveToMemory(ctx context.Context, data *dto.UserIdInstance) ([]int, error) {
	var numbers []int
	var err error

	if numbers, err = r.persistent.InstancePermissionsNumbersForAccount(ctx, data); err == nil && len(numbers) > 0 {
		go func() {
			_ = r.memory.SetInstancePermissionsNumbersForAccount(context.Background(), &dto.UserIdInstancePermNumbers{
				UserId:            data.UserId,
				Instance:          data.Instance,
				PermissionNumbers: numbers,
			})
		}()
	}

	return numbers, adaptErr(err)
}

// saveToMemoryLoginData сохраняет в памяти данные, необходимые для процесса входа в систему пользователя (сервиса).
func (r *Repository) saveToMemoryLoginData(ctx context.Context, data *dto.UserIdLoginHashState) error {
	r.memory.SetUserIdAndPasswordHash(ctx, &dto.UserIdLoginHash{UserId: data.UserId, Login: data.Login, Hash: data.Hash})
	return adaptErr(r.memory.SetAccountState(ctx, &dto.LoginState{Login: data.Login, State: data.State}))
}

// refreshAccountPermissions обновляет кеш разрешений
func (r *Repository) refreshAccountPermissions(ctx context.Context, data *dto.UserIdService) {
	if servicePerm, err := r.persistent.ServicePermissionsNumbersForAccount(ctx, data); err == nil {
		_ = r.memory.SetServicePermissionsNumbersForAccount(context.Background(), &dto.UserIdServicePermNumbers{
			UserId:            data.UserId,
			Service:           data.Service,
			PermissionNumbers: servicePerm,
		})
	}
}

// makeDataCache считывает все данные (которые возможно кешировать) из постоянного хранилища в хранилище в памяти.
func (r *Repository) makeDataCache() {
	slog.Info("data caching has started")
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		slog.Info(fmt.Sprintf("data caching is complete (time spent on caching %s)", elapsed.String()))
	}()

	ctx := context.Background()
	enabledAccountsLogins, err := r.persistent.AccountsLoginsByState(ctx, account_state.Enabled)
	if err != nil {
		slog.Error(err.Error())
	} else {
		c := make(chan struct{}, r.persistent.MaxConnections()/2)
		defer close(c)

		for _, login := range enabledAccountsLogins {
			c <- struct{}{}

			go func(login loginVO.Login) {
				if _, err = r.AccountLoginData(ctx, login); err != nil {
					slog.Error(err.Error())
				}
			}(login)

			<-c
		}
	}

	var services []string

	if services, err = r.ServicesNames(ctx); err == nil {
		for _, service := range services {
			_, _ = r.ServiceNumberedPermissions(ctx, service)
		}
	}
}
