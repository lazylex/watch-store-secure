package lexerr

const persistentType = "persistent repo"

var (
	ErrDuplicateKeyValue = NewPersistentError("duplicate key value violates unique constraint violation")
	ErrZeroRowsAffected  = NewPersistentError("zero rows affected")
)

type Persistent struct {
	BaseError
}

// FullPersistentError возвращает полностью заполненную структуру Persistent
func FullPersistentError(message, origin string, initialError error) *Persistent {
	p := &Persistent{}
	p.Type = persistentType
	p.Message = message
	p.Origin = origin
	p.InitialError = initialError

	return p
}

// NewPersistentError возвращает структуру ошибки Persistent с переданным в качестве аргумента сообщением
func NewPersistentError(message string) *Persistent {
	p := &Persistent{}
	p.Type = persistentType
	p.Message = message

	return p
}

// WithOrigin добавляет в структуру место появления ошибки
func (p *Persistent) WithOrigin(origin string) *Persistent {
	p.Origin = origin
	return p
}
