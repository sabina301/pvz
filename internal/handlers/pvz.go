package handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"pvz/internal/bootstrap"
	"pvz/internal/logger"
	"pvz/internal/models/pvz"
	"pvz/internal/services"
	"pvz/pkg/errors"
	"strconv"
	"time"
)

type PvzHandler struct {
	pvzService services.PvzService
}

func NewPvzHandler(deps bootstrap.Deps) *PvzHandler {
	return &PvzHandler{
		pvzService: deps.PvzService,
	}
}

func (ph *PvzHandler) Create(c echo.Context) error {
	log := logger.Log.With("handler", "pvz", "method", "Create")
	log.Info("received request to create pvz")

	var req pvz.CreateRequest
	if err := c.Bind(&req); err != nil {
		log.Error("failed to bind request body", "error", err)
		return errors.NewMalformedBody()
	}

	if err := apiValidator.ValidateRequest(req); err != nil {
		log.Error("request validation failed", "error", err)
		return err
	}

	createdPvz, err := ph.pvzService.Create(req)
	if err != nil {
		log.Error("pvzService.Create failed", "error", err)
		return err
	}
	log.Info("pvz successfully created", "pvzId", createdPvz.Id)

	return c.JSON(http.StatusCreated, createdPvz)
}

func (ph *PvzHandler) DeleteLastProduct(c echo.Context) error {
	log := logger.Log.With("handler", "pvz", "method", "DeleteLastProduct")

	pvzId := c.Param("pvzId")
	log.Info("received request to delete last product", "pvzId", pvzId)

	if err := apiValidator.ValidateParam(pvzId, "required,uuid"); err != nil {
		log.Error("parameter validation failed", "error", err)
		return err
	}

	updatedPvz, err := ph.pvzService.DeleteLastProduct(pvzId)
	if err != nil {
		log.Error("pvzService.DeleteLastProduct failed", "error", err)
		return err
	}
	log.Info("last product successfully deleted", "pvzId", updatedPvz.Id)

	return c.JSON(http.StatusOK, updatedPvz)
}

func (ph *PvzHandler) CloseLastReception(c echo.Context) error {
	log := logger.Log.With("handler", "pvz", "method", "CloseLastReception")

	pvzId := c.Param("pvzId")
	log.Info("received request to close last reception", "pvzId", pvzId)

	if err := apiValidator.ValidateParam(pvzId, "required,uuid"); err != nil {
		log.Error("parameter validation failed", "error", err)
		return err
	}

	updatedPvz, err := ph.pvzService.CLoseLastReception(pvzId)
	if err != nil {
		log.Error("pvzService.CLoseLastReception failed", "error", err)
		return err
	}
	log.Info("last reception successfully closed", "pvzId", updatedPvz.Id)

	return c.JSON(http.StatusOK, updatedPvz)
}

func (ph *PvzHandler) ListWithFilterDate(c echo.Context) error {
	log := logger.Log.With("handler", "pvz", "method", "ListWithFilterDate")

	startDateStr := c.QueryParam("startDate")

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return errors.NewBadParamValue(startDateStr)
	}
	
	endDateStr := c.QueryParam("endDate")
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return errors.NewBadParamValue(endDateStr)
	}
	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return errors.NewBadParamValue(page)
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return errors.NewBadParamValue(limit)
	}

	req := pvz.ListRequest{
		StartDate: &startDate,
		EndDate:   &endDate,
		Page:      pageInt,
		Limit:     limitInt,
	}
	list, err := ph.pvzService.ListWithFilterDate(req)
	if err != nil {
		log.Error("pvzService.ListWithFilterDate failed", "error", err)
		return err
	}
	log.Info("received pvz list with filter date", "count", len(list))

	return c.JSON(http.StatusOK, list)
}
