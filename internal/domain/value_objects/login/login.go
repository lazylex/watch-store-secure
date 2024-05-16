package login

import (
	"fmt"
)

type Login string

const (
	minLength = 3
	maxLength = 100
)

func (l Login) Validate() error {
	if len(l) > maxLength {
		return fmt.Errorf("maximum login length exceeded (%d characters)", maxLength)
	}
	if len(l) < minLength {
		return fmt.Errorf("login length must be at least %d characters", minLength)
	}

	return nil
}
