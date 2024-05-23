package handlers

import (
	"context"
	"fmt"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/service"
	"log/slog"
	"net/http"
	"time"
)

// Handler структура для обработки http-запросов.
type Handler struct {
	service      *service.Service // Объект, реализующий логику сервиса
	queryTimeout time.Duration    // Допустимый таймаут для обработки запроса
}

// New возвращает структуру с обработчиками http-запросов.
func New(domainService *service.Service, timeout time.Duration) *Handler {
	return &Handler{service: domainService, queryTimeout: timeout}
}

// Login производит вход в учетную запись и возвращает в JSON токен сессии(по ключу token).Тип авторизации - Basic Auth.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodPost, w, r) {
		return
	}

	var ok bool
	var err error
	var token, username, pwd string

	if username, pwd, ok = r.BasicAuth(); !ok {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"secure\"")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userLogin := login.Login(username)
	userPassword := password.Password(pwd)

	if userLogin.Validate() != nil || userPassword.Validate() != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if token, err = h.service.Login(ctx, &dto.LoginPassword{Login: userLogin, Password: userPassword}); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", token)))
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Here")
}

// allowedOnlyMethod принимает разрешенный метод и, если запрос ему не соответствует, записывает в заголовок информацию
// о разрешенном методе, статус http.StatusMethodNotAllowed и возвращает false.
func allowedOnlyMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}

	return true
}
