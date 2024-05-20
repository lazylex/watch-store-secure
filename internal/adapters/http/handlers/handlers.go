package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/service"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	service      *service.Service
	queryTimeout time.Duration
}

func New(domainService *service.Service, timeout time.Duration) *Handler {
	return &Handler{service: domainService, queryTimeout: timeout}
}

// Login производит вход в учетную запись и возвращает в JSON токен сессии(по ключу token).Тип авторизации - Basic Auth.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !allowedOnlyMethod(http.MethodPost, w, r) {
		return
	}

	const prefix = "Basic "

	var (
		decodedBytes []byte
		err          error
		token        string
	)

	auth := r.Header.Get("Authorization")
	if len(auth) < len(prefix) {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"secure\"")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if decodedBytes, err = base64.StdEncoding.DecodeString(auth[len(prefix):]); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authData := strings.Split(string(decodedBytes), ":")
	if len(authData) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userLogin := login.Login(authData[0])
	pwd := password.Password(authData[1])

	if userLogin.Validate() != nil || pwd.Validate() != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.queryTimeout)
	defer cancel()

	if token, err = h.service.Login(ctx, &dto.LoginPassword{Login: userLogin, Password: pwd}); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", token)))
}

func allowedOnlyMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}

	return true
}
