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
				log := slog.Default().With("origin", "recoverer middleware")

				switch t := rvr.(type) {
				case error:
					log.Warn(t.Error())
				case string:
					log.Warn(t)
				default:
					log.Warn(fmt.Sprint(rvr))
				}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
