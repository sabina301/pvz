package internal

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"pvz/configs"
	"pvz/internal/bootstrap"
)

type Server struct {
	e   *echo.Echo
	cfg *configs.AppConfig
}

func NewServer(cfg *configs.AppConfig, deps bootstrap.Deps) *Server {
	e := echo.New()

	createApi(e, deps)
	return &Server{
		e:   e,
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	return s.e.Start(fmt.Sprintf(":%d", s.cfg.HttpPort))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
