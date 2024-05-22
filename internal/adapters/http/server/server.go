package server

import (
	"context"
	"errors"
	"github.com/lazylex/watch-store/secure/internal/adapters/http/handlers"
	"github.com/lazylex/watch-store/secure/internal/adapters/http/middleware/token_checker"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/service"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Server struct {
	cfg     *config.HttpServer
	paths   []string
	srv     *http.Server
	mux     *http.ServeMux
	service *service.Service
}

func MustCreate(domainService *service.Service, cfg *config.HttpServer) *Server {
	mux := http.NewServeMux()
	server := &Server{mux: mux, service: domainService, cfg: cfg}
	server.srv = &http.Server{
		Addr:         server.cfg.Address,
		Handler:      server.mux,
		ReadTimeout:  server.cfg.ReadTimeout,
		WriteTimeout: server.cfg.WriteTimeout,
		IdleTimeout:  server.cfg.IdleTimeout,
	}
	// TODO заменить задержку с 10 секунд на чтение из конфигурации
	h := handlers.New(domainService, 10*time.Second)
	server.assignPathToHandler("/login", h.Login)

	tokenMiddleware := token_checker.New(domainService)
	server.srv.Handler = tokenMiddleware.Checker(server.mux)

	return server
}

// assignPathToHandler проверяет, не прикреплен ли уже переданный первым аргументом функции адрес к какому-либо
// обработчику. Если прикреплен, то выполнение функции прекращается, чтобы не вызвать панику в http.HandleFunc. При
// нормальном выполнении, добавляет пусть к списку используемых и прикрепляет его к переданному вторым аргументом
// обработчику.
func (s *Server) assignPathToHandler(path string, handler func(http.ResponseWriter, *http.Request)) {
	for _, v := range s.paths {
		if v == path {
			return
		}
	}

	s.paths = append(s.paths, path)
	s.mux.HandleFunc(path, handler)
}

// MustRun производит запуск сервера в отдельной go-рутине. В случае ошибки останавливает работу приложения.
func (s *Server) MustRun() {
	go func() {
		slog.Info("start http server on " + s.srv.Addr)
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server err: startup error. Initial error: " + err.Error())
			os.Exit(1)
		}
	}()
}

// Shutdown производит остановку сервера.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("failed to gracefully shutdown http server")
	} else {
		slog.Info("gracefully shut down http server")
	}
}
