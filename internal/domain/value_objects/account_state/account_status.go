package account_state

type State int

const (
	Enabled = iota + 1
	Disabled
)

// IsStateCorrect возвращает true, если переданное состояние определено в программе
func IsStateCorrect(state State) bool {
	if state == Enabled || state == Disabled {
		return true
	}
	return false
}
