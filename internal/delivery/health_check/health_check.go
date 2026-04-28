package health_check

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	redis *redis.Client
	db    *sql.DB
	l     *slog.Logger
}

func NewHealthCheck(redis *redis.Client, db *sql.DB, l *slog.Logger) *HealthHandler {
	return &HealthHandler{
		redis: redis,
		db:    db,
		l:     l,
	}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {

	err := h.db.PingContext(r.Context())
	if err != nil {
		h.l.Error("ошибка подключения к db", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.redis.Ping(r.Context()).Err()
	if err != nil {
		h.l.Error("ошибка подключения к redis", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

func (h *HealthHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong\n"))
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ping", h.Ping)
}
