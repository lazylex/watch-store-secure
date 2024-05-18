/*
Package errors: пакет для определения структурированной ошибки.
*/
package errors

import (
	"fmt"
	"runtime"
)

// BaseError структура базовой ошибки, на основе которой заполняются специализированные ошибки. Реализует интерфейс
// error.
type BaseError struct {
	Type         string // Тип ошибки
	Message      string // Сообщение ошибки
	Origin       string // Место возникновения ошибки
	InitialError error  // Ошибка, повлекшая возникновение текущей ошибки
}

// Error возвращает текстовое описание ошибки.
func (b *BaseError) Error() string {
	var result string

	if len(b.Type) > 0 {
		result += b.Type + " err"
	}

	if len(b.Message) > 0 {
		result = fmt.Sprintf("%s: %s.", result, b.Message)
	} else {
		result = fmt.Sprintf("%s.", result)
	}

	if b.InitialError != nil {
		result = fmt.Sprintf("%s Initial err: %s.", result, b.InitialError.Error())
	}

	if len(b.Origin) > 0 {
		result = fmt.Sprintf("%s Origin: %s", result, b.Origin)
	}

	return result

}

// GetFrame возвращает фрейм на переданной глубине вложенности.
func GetFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// WithOrigin сохраняет в структуре место появления ошибки.
func (b *BaseError) WithOrigin(origin string) *BaseError {
	b.Origin = origin
	return b
}
