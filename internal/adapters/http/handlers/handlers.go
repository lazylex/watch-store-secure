package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	serviceErr "github.com/lazylex/watch-store/secure/internal/errors/service"
	v "github.com/lazylex/watch-store/secure/internal/helpers/constants/various"
	"github.com/lazylex/watch-store/secure/internal/service"
	"net/http"
	"time"
)

// Проверка на существование токена содержится в middleware, поэтому в обработчиках она опускается.

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
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusRequestTimeout)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", token)))
}

// Index обработчик для несуществующих страниц.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
}

// Logout производит выход из учетной записи.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodGet, w, r) {
		return
	}

	token := r.Header.Get("Authorization")[len(v.BearerTokenPrefix):]

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	id, err := h.service.GetUserUUIDFromSession(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.service.Logout(ctx, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetTokenWithPermissions возвращает JWT-токен, содержащий информацию о разрешениях для переданного экземпляра
// приложения.
func (h *Handler) GetTokenWithPermissions(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodGet, w, r) {
		return
	}

	var (
		err   error
		token string
		id    uuid.UUID
	)

	instance := r.FormValue("instance")

	if len(instance) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token = r.Header.Get("Authorization")[len(v.BearerTokenPrefix):]

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if id, err = h.service.GetUserUUIDFromSession(ctx, token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if token, err = h.service.CreateToken(ctx, &dto.UserIdInstance{UserId: id, Instance: instance}); err != nil {
		if errors.Is(err, serviceErr.ErrEmptyResult) {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"jwt-token\":\"%s\"}", token)))
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
