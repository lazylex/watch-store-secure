package in_memory

import "github.com/lazylex/watch-store/secure/internal/dto"

type Interface interface {
	SaveSession(dto.SessionDTO) error
}
