package joint

import (
	"context"
	voLogin "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/in_memory"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/persistent"
)

type Repository struct {
	memory     in_memory.Interface
	persistent persistent.Interface
}

func New(memory in_memory.Interface, persistent persistent.Interface) *Repository {
	return &Repository{memory: memory, persistent: persistent}
}

func (r *Repository) SaveSession(ctx context.Context, dto dto.SessionDTO) error {
	return r.memory.SaveSession(dto)
}

func (r *Repository) SaveLoginAndPasswordHash(ctx context.Context, dto dto.LoginWithPasswordHashDTO) error {
	// TODO implement
	return nil
}

func (r *Repository) GetIdAndPasswordHash(ctx context.Context, login voLogin.Login) (dto.IdWithPasswordHashDTO, error) {
	// TODO implement
	return dto.IdWithPasswordHashDTO{}, nil
}
