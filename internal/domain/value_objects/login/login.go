package login

import (
	"fmt"
)

type Login string

const (
	minLength = 3
	maxLength = 100
)

// Validate возвращает ошибку, если логин некорректный.
func (l Login) Validate() error {
	if len(l) > maxLength {
		return fmt.Errorf("maximum login length exceeded (%d characters)", maxLength)
	}
	if len(l) < minLength {
		return fmt.Errorf("login length must be at least %d characters", minLength)
	}

	// RFC 7617 'Basic' HTTP Authentication Scheme September 2015 запрещает использование двоеточия в логине
	for _, c := range l {
		if c == ':' {
			return fmt.Errorf("login should not contain ':' character")
		}
	}

	return nil
}
