package message_broker

import "github.com/lazylex/watch-store/secure/internal/errors"

const brokerType = "message broker"

var (
	ErrCouldNotSendMessage = NewMessageBrokerError("couldn't send message after several attempts")
	ErrFailedToCloseWriter = NewMessageBrokerError("failed to close writer")
)

// FullMessageBrokerError возвращает полностью заполненную структуру с типом MessageBrokerType.
func FullMessageBrokerError(message, origin string, initialError error) *errors.BaseError {
	return &errors.BaseError{
		Type:         brokerType,
		Message:      message,
		Origin:       origin,
		InitialError: initialError,
	}
}

// NewMessageBrokerError возвращает структуру ошибки с типом MessageBrokerType и переданным в качестве аргумента
// сообщением.
func NewMessageBrokerError(message string) *errors.BaseError {
	return &errors.BaseError{Type: brokerType, Message: message}
}
