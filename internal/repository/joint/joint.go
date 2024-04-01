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

// saveToMemoryLoginData сохраняет в памяти данные, необходимые для процесса входа в систему пользователя (сервиса)
func (r *Repository) saveToMemoryLoginData(ctx context.Context, data dto.AccountLoginDataDTO) error {
	r.memory.SetUserIdAndPasswordHash(ctx, dto.UserLoginAndIdWithPasswordHashDTO{UserId: data.UserId, Login: data.Login, Hash: data.Hash})
	return r.memory.SetAccountState(ctx, dto.LoginStateDTO{Login: data.Login, State: data.State})
}

// makeDataCache считывает все данные (которые возможно кешировать) из постоянного хранилища в хралищие в памяти
func makeDataCache() {
	// TODO implement
	slog.Debug(jointRepositoryError("makeDataCache not implemented").Error())
}
