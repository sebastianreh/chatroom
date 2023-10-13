package ping

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"

	"github.com/sebastianreh/chatroom/internal/config"
)

type Handler interface {
	Ping(c echo.Context) error
}

type Response struct {
	Version string    `json:"version"`
	Name    string    `json:"name"`
	Uptime  time.Time `json:"uptime"`
}

type StatusHandler struct {
	config config.Config
}

func NewSHandierPing(cfg config.Config) Handler {
	return &StatusHandler{
		config: cfg,
	}
}

func (s *StatusHandler) Ping(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Response{
		Version: s.config.ProjectVersion,
		Name:    s.config.ProjectName,
		Uptime:  time.Now().UTC(),
	})
}
