package joint

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/in_memory"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/persistent"
)

type Repository struct {
	memory     in_memory.Interface
	persistent persistent.Interface
}

func New(memory in_memory.Interface, persistent persistent.Interface) Repository {
	go makeDataCache()
	return Repository{memory: memory, persistent: persistent}
}

func (r *Repository) SaveSession(ctx context.Context, dto dto.SessionDTO) error {
	return r.memory.SaveSession(dto)
}

func (r *Repository) SaveLoginAndPasswordHash(ctx context.Context, dto dto.LoginWithPasswordHashDTO) error {
	// TODO implement
	return nil
}

func (r *Repository) GetUserIdAndPasswordHash(ctx context.Context, login login.Login) (dto.UserIdWithPasswordHashDTO, error) {
	// TODO implement
	return dto.UserIdWithPasswordHashDTO{}, nil
}

func (r *Repository) GetId(ctx context.Context, login login.Login) (uuid.UUID, error) {
	// TODO implement
	return uuid.UUID{}, nil
}

// makeDataCache считывает все данные (которые возможно кешировать) из постоянного хранилища в хралищие в памяти
func makeDataCache() {
	// TODO implement
}
