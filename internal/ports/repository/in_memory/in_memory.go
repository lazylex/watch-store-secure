package in_memory

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type Interface interface {
	SaveSession(context.Context, dto.SessionDTO) error
}
