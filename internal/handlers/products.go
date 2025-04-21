package handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"pvz/internal/bootstrap"
	"pvz/internal/logger"
	"pvz/internal/models/product"
	"pvz/internal/services"
	"pvz/pkg/errors"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(deps bootstrap.Deps) *ProductHandler {
	return &ProductHandler{
		productService: deps.ProductService,
	}
}

func (ph *ProductHandler) AddInReception(c echo.Context) error {
	log := logger.Log.With("handler", "product", "method", "AddInReception")

	var req product.AddInReceptionRequest
	if err := c.Bind(&req); err != nil {
		log.Error("failed to bind request body", "error", err)
		return errors.NewMalformedBody()
	}

	if err := apiValidator.ValidateRequest(req); err != nil {
		log.Error("request validation failed", "error", err)
		return err
	}

	log.Info("calling productService.AddInReception", "pvzId", req.PvzId)
	createdProduct, err := ph.productService.AddInReception(req)
	if err != nil {
		log.Error("productService.AddInReception failed", "error", err)
		return err
	}
	log.Info("product added to reception successfully", "productId", createdProduct.Id)

	return c.JSON(http.StatusCreated, createdProduct)
}
