package v1

import (
	"context"
	"debez/internal/repository"
	"debez/internal/service"
	"debez/internal/transport/http/handlers"
	"debez/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	srv *http.Server
	db  *pgxpool.Pool
}

const (
	defaultHeaderTimeout = 5 * time.Second
)

func NewServer(port int, db *pgxpool.Pool) *Server {
	s := http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           nil,
		ReadHeaderTimeout: defaultHeaderTimeout,
	}
	return &Server{
		srv: &s,
		db:  db,
	}
}
func (s *Server) RegisterHandler(ctx context.Context) error {
	userRepo := repository.NewUserRepository(s.db)
	userService := service.NewUserService(userRepo)
	handler := handlers.NewHandlerFacade(ctx, userService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/users", logger.LoggerMiddleware(ctx, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetUsers(w, r)
	}))
	mux.HandleFunc("/api/v1/users/{id}", logger.LoggerMiddleware(ctx, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.GetUserByID(w, r)
	}))
	mux.HandleFunc("/api/v1/create_user", logger.LoggerMiddleware(ctx, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.SaveUser(w, r)
	}))
	mux.HandleFunc("/api/v1/update_user", logger.LoggerMiddleware(ctx, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.UpdateUser(w, r)
	}))
	mux.HandleFunc("/api/v1/delete_user/{id}", logger.LoggerMiddleware(ctx, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.DeleteUser(w, r)
	}))
	s.srv.Handler = LoggingMiddleware(ctx)(mux)

	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}
func (s *Server) Stop(ctx context.Context) error {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "shutting down http server")
	return s.srv.Shutdown(ctx)
}
