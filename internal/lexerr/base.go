package lexerr

import (
	"fmt"
	"runtime"
)

type BaseError struct {
	Type         string
	Message      string
	Origin       string
	InitialError error
}

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
