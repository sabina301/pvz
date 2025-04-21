package internal

import (
	"github.com/labstack/echo"
	"pvz/internal/bootstrap"
	"pvz/internal/handlers"
	"pvz/internal/middleware"
	aModel "pvz/internal/models/auth"
)

const (
	serviceVersion = "v1"
	apiPrefix      = "/api/" + serviceVersion
)

func createApi(e *echo.Echo, deps bootstrap.Deps) {
	api := e.Group(apiPrefix, middleware.SetApiTimeout, middleware.HandleError)

	authHandler := handlers.NewAuthHandler(deps)
	auth := api.Group("", middleware.SetApiTimeout)
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/dummyLogin", authHandler.DummyLogin)

	pvzHandler := handlers.NewPvzHandler(deps)
	pvz := api.Group("/pvz", middleware.SetApiTimeout)
	pvz.POST("", pvzHandler.Create, middleware.AllowRoles(aModel.Moderator))
	pvz.POST("/:pvzId/delete_last_product", pvzHandler.DeleteLastProduct, middleware.AllowRoles(aModel.Employee))
	pvz.POST("/:pvzId/close_last_reception", pvzHandler.CloseLastReception, middleware.AllowRoles(aModel.Employee))
	pvz.GET("", pvzHandler.ListWithFilterDate)

	receptionHandler := handlers.NewReceptionHandler(deps)
	reception := api.Group("/receptions", middleware.SetApiTimeout)
	reception.POST("", receptionHandler.Create, middleware.AllowRoles(aModel.Employee))

	productHandler := handlers.NewProductHandler(deps)
	product := api.Group("/product", middleware.SetApiTimeout)
	product.POST("", productHandler.AddInReception, middleware.AllowRoles(aModel.Employee))
}
