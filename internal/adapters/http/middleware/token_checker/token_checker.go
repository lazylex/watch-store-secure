package token_checker

import (
	v "github.com/lazylex/watch-store/secure/internal/helpers/constants/various"
	"github.com/lazylex/watch-store/secure/internal/service"
	"net/http"
	"strings"
)

const loginURI = "/login"

// TokenChecker структура, содержащая доступ к сервисной логике.
type TokenChecker struct {
	service *service.Service
}

// New служит для создания middleware, предназначенного для отклонения запросов, не содержащих токен или не
// предназначенных для входа в систему.
func New(service *service.Service) *TokenChecker {
	return &TokenChecker{service: service}
}

// Checker проверяет, что запрос либо осуществляется по адресу, назначенному для процедуры входа в систему, либо
// содержит токен, который соответствует открытой сессии.
func (t *TokenChecker) Checker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == loginURI {
			next.ServeHTTP(w, req)
			return
		}

		authHeader := req.Header.Get("Authorization")

		if len(authHeader) == 0 || !strings.HasPrefix(authHeader, v.BearerTokenPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if _, err := t.service.UserUUIDFromSession(req.Context(), authHeader[len(v.BearerTokenPrefix):]); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}
