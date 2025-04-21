package handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"pvz/internal/bootstrap"
	"pvz/internal/logger"
	"pvz/internal/models/auth"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

func NewAuthHandler(deps bootstrap.Deps) *AuthHandler {
	return &AuthHandler{
		authService: deps.AuthService,
	}
}

type AuthHandler struct {
	authService services.AuthService
}

func (ah *AuthHandler) Login(c echo.Context) error {
	log := logger.Log.With("handler", "auth", "method", "Login")

	var req auth.LoginRequest
	if err := c.Bind(&req); err != nil {
		log.Error("failed to bind request body", "error", err)
		return errors.NewMalformedBody()
	}

	if err := apiValidator.ValidateRequest(req); err != nil {
		log.Error("request validation failed", "error", err)
		return err
	}

	log.Info("calling authService.Login")
	token, err := ah.authService.Login(req)
	if err != nil {
		log.Error("authService.Login failed", "error", err)
		return err
	}

	log.Info("login successful")
	return c.JSON(http.StatusOK, token)
}

func (ah *AuthHandler) Register(c echo.Context) error {
	log := logger.Log.With("handler", "auth", "method", "Register")

	var req auth.RegisterRequest
	if err := c.Bind(&req); err != nil {
		log.Error("failed to bind request body", "error", err)
		return errors.NewMalformedBody()
	}

	log.Info("validating request", "req", req)
	if err := apiValidator.ValidateRequest(req); err != nil {
		log.Error("request validation failed", "error", err)
		return err
	}

	log.Info("calling authService.Register")
	user, err := ah.authService.Register(req)
	if err != nil {
		log.Error("authService.Register failed", "error", err)
		return err
	}

	log.Info("registration successful")
	return c.JSON(http.StatusOK, user)
}

func (ah *AuthHandler) DummyLogin(c echo.Context) error {
	log := logger.Log.With("handler", "auth", "method", "DummyLogin")

	var req auth.DummyLoginRequest
	if err := c.Bind(&req); err != nil {
		log.Error("failed to bind request body", "error", err)
		return errors.NewMalformedBody()
	}

	if err := apiValidator.ValidateRequest(req); err != nil {
		log.Error("request validation failed", "error", err)
		return err
	}

	log.Info("calling authService.DummyLogin")
	token, err := ah.authService.DummyLogin(req)
	if err != nil {
		log.Error("authService.DummyLogin failed", "error", err)
		return err
	}

	log.Info("dummy login successful")
	return c.JSON(http.StatusOK, token)
}
