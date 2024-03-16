package joint

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/in_memory"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/persistent"
	"log/slog"
)

type Repository struct {
	memory     in_memory.Interface
	persistent persistent.Interface
}

var (
	ErrNoRecord = jointRepositoryError("no record")
)

func jointRepositoryError(text string) error {
	return errors.New("joint repo: " + text)
}

func New(memory in_memory.Interface, persistent persistent.Interface) Repository {
	go makeDataCache()
	return Repository{memory: memory, persistent: persistent}
}

func (r *Repository) SaveSession(ctx context.Context, dto dto.SessionDTO) error {
	return r.memory.SaveSession(ctx, dto)
}

func (r *Repository) SaveLoginAndPasswordHash(ctx context.Context, dto dto.LoginWithPasswordHashDTO) error {
	// TODO implement
	slog.Debug(jointRepositoryError("SaveLoginAndPasswordHash not implemented").Error())
	return nil
}

// GetUserIdAndPasswordHash по переданному логину возвращает идентификатор пользователя и хэш его пароля. Поиск
// производится сперва в in memory хранилище, затем в постоянной БД. Если в памяти нет, а в постоянной БД есть, то
// осуществляется попытка кеширования данных в память
func (r *Repository) GetUserIdAndPasswordHash(ctx context.Context, login login.Login) (dto.UserIdWithPasswordHashDTO, error) {
	var data dto.UserIdWithPasswordHashDTO
	var err error

	if data, err = r.memory.GetUserIdAndPasswordHash(ctx, login); err != nil {
		///////////////////////////////

		if data, err = r.persistent.GetUserIdAndPasswordHash(ctx, login); err != nil {
			return dto.UserIdWithPasswordHashDTO{}, err
		} else if data != (dto.UserIdWithPasswordHashDTO{}) {
			// если не смогли кешировать в редис запись, ничего страшного
			defer r.memory.SetUserIdAndPasswordHash(ctx,
				dto.UserLoginAndIdWithPasswordHashDTO{Login: login, UserId: data.UserId, Hash: data.Hash})
		} else {
			return dto.UserIdWithPasswordHashDTO{}, ErrNoRecord
		}
	}

	return data, nil
}

func (r *Repository) GetId(ctx context.Context, login login.Login) (uuid.UUID, error) {
	// TODO implement
	slog.Debug(jointRepositoryError("GetId not implemented").Error())
	return uuid.UUID{}, nil
}

// makeDataCache считывает все данные (которые возможно кешировать) из постоянного хранилища в хралищие в памяти
func makeDataCache() {
	// TODO implement
	slog.Debug(jointRepositoryError("makeDataCache not implemented").Error())
}
