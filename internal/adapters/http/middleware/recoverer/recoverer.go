package recoverer

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Recoverer middleware который восстанавливает после паники и заносит данные о причине в лог. Так же возвращает HTTP
// статус 500 (Internal Server Error) если это возможно.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				writeToLog(rvr)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// writeToLog записывает в лог причину паники.
func writeToLog(rvr any) {
	log := slog.Default().With("origin", "recovery middleware")

	switch t := rvr.(type) {
	case error:
		log.Warn("panic error: " + t.Error())
	case string:
		log.Warn("panic string: " + t)
	default:
		log.Warn(fmt.Sprint(rvr))
	}
}
