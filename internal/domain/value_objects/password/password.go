package password

import (
	"fmt"
	"strings"
)

type Password string

const minLength = 8

// Validate возвращает ошибку, если пароль некорректный.
func (p Password) Validate() error {
	if len(p) < minLength {
		return fmt.Errorf("password length must be at least %d characters", minLength)
	}

	if strings.ToUpper(string(p)) == string(p) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if strings.ToLower(string(p)) == string(p) {
		return fmt.Errorf("password must contain at least ont uppercase letter")
	}

	return nil
}
