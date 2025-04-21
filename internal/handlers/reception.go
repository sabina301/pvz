package handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"pvz/internal/bootstrap"
	"pvz/internal/logger"
	"pvz/internal/models/reception"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

type ReceptionHandler struct {
	receptionService services.ReceptionService
}

func NewReceptionHandler(deps bootstrap.Deps) *ReceptionHandler {
	return &ReceptionHandler{
		receptionService: deps.ReceptionService,
	}
}

func (rh *ReceptionHandler) Create(c echo.Context) error {
	log := logger.Log.With("handler", "reception", "method", "Create")

	var req reception.CreateRequest
	if err := c.Bind(&req); err != nil {
		log.Error("failed to bind request body", "error", err)
		return errors.NewMalformedBody()
	}

	if err := apiValidator.ValidateRequest(req); err != nil {
		log.Error("request validation failed", "error", err)
		return err
	}

	newReception, err := rh.receptionService.Create(req)
	if err != nil {
		log.Error("receptionService.Create failed", "error", err)
		return err
	}
	log.Info("reception successfully created", "receptionId", newReception.Id, "pvzId", newReception.PvzId)

	return c.JSON(http.StatusCreated, newReception)
}
