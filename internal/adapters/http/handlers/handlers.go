package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	serviceErr "github.com/lazylex/watch-store/secure/internal/errors/service"
	v "github.com/lazylex/watch-store/secure/internal/helpers/constants/various"
	"github.com/lazylex/watch-store/secure/internal/service"
	"log/slog"
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

// Login производит вход в учетную запись и возвращает в JSON токен сессии(по ключу token). Тип авторизации - Basic Auth.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodPost, w, r) {
		return
	}

	var ok bool
	var err error
	var token, username, pwd string
	var log = slog.Default().With("remote address", r.RemoteAddr)

	if username, pwd, ok = r.BasicAuth(); !ok {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"secure\"")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userLogin := login.Login(username)
	userPassword := password.Password(pwd)

	if userLogin.Validate() != nil || userPassword.Validate() != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Warn("unable to validate username or password")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if token, err = h.service.Login(ctx, &dto.LoginPassword{Login: userLogin, Password: userPassword}); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			w.WriteHeader(http.StatusRequestTimeout)
			log.Warn("request timed out")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			log.Warn("unable to login")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", token)))

	log.Info("successfully logged in")
}

// Index обработчик для несуществующих страниц.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		slog.Default().With("remote address", r.RemoteAddr).With("request url", r.RequestURI).Warn("page not found")
		return
	}
}

// Logout производит выход из учетной записи.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodGet, w, r) {
		return
	}

	log := slog.Default().With("remote address", r.RemoteAddr)
	token := r.Header.Get("Authorization")[len(v.BearerTokenPrefix):]

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	id, err := h.service.UserUUIDFromSession(ctx, token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("unable to get user uuid from session")
		return
	}

	if err = h.service.Logout(ctx, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("unable to logout")
		return
	}

	log.Info("successfully logout")
}

// TokenWithPermissions возвращает JWT-токен, содержащий информацию о разрешениях для переданного экземпляра приложения.
func (h *Handler) TokenWithPermissions(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodGet, w, r) {
		return
	}

	var (
		err   error
		token string
		id    uuid.UUID
		log   = slog.Default().With("remote address", r.RemoteAddr)
	)

	instance := r.FormValue("instance")

	if len(instance) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn("unable to get instance")
		return
	}

	token = r.Header.Get("Authorization")[len(v.BearerTokenPrefix):]

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if id, err = h.service.UserUUIDFromSession(ctx, token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Warn("unable to get user uuid from session")
		return
	}

	if token, err = h.service.CreateToken(ctx, &dto.UserIdInstance{UserId: id, Instance: instance}); err != nil {
		if errors.Is(err, serviceErr.ErrEmptyResult) {
			w.WriteHeader(http.StatusNoContent)
			slog.Warn("no token for return")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Warn("error create token")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"jwt-token\":\"%s\"}", token)))

	log.Info("sent jwt-token")
}

// ServiceNumberedPermissions возвращает JSON с названиями и номерами разрешений для переданного в параметре service
// сервиса. При отсутствии разрешений возвращает статус http.StatusNoContent.
// Пример возвращаемого функцией JSON:
//
// [
//
//	{
//		"name": "удалять данные",
//		"number": 1
//	},
//	{
//		"name": "добавлять данные",
//		"number": 5
//	},
//
// ]
func (h *Handler) ServiceNumberedPermissions(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodGet, w, r) {
		return
	}

	var (
		answer              []byte
		numberedPermissions *[]dto.NameNumber
		err                 error
		log                 = slog.Default().With("remote address", r.RemoteAddr)
	)

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	serviceName := r.FormValue("service")

	if len(serviceName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Warn("unable to get service")
	}

	numberedPermissions, err = h.service.ServiceNumberedPermissions(ctx, serviceName)
	if err != nil {
		if errors.Is(err, serviceErr.ErrEmptyResult) {
			w.WriteHeader(http.StatusNoContent)
			log.Warn("no service numbered permissions")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Warn("error get numbered permissions")
		}
		return
	}

	if answer, err = json.Marshal(numberedPermissions); err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(answer)
		log.Info("numbered service permits have been sent")

	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Warn("unable to marshal numbered permissions")
	}
}

// allowedOnlyMethod принимает разрешенный метод и, если запрос ему не соответствует, записывает в заголовок информацию
// о разрешенном методе, статус http.StatusMethodNotAllowed и возвращает false.
func allowedOnlyMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Default().With("remote address", r.RemoteAddr).With("request url", r.RequestURI).Warn("method not allowed")
		return false
	}

	return true
}
