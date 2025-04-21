package middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"pvz/configs"
	"pvz/internal/logger"
	"pvz/internal/models/auth"
	"pvz/internal/tokens"
	"pvz/pkg/errors"
	"slices"
	"strings"
)

func SetApiTimeout(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		timeout := configs.DefaultAPIRequestTimeout
		ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
		defer cancel()
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func AllowRoles(trueRoles ...auth.Role) func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !HasAccess(c.Request(), trueRoles) {
				logger.Log.Error("access denied")
				return errors.NewAccessForbidden()
			}
			return next(c)
		}
	}
}

func HasAccess(r *http.Request, roles []auth.Role) bool {
	header := r.Header.Get("Authorization")
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return false
	}
	jwt := parts[1]
	claims, err := tokens.ParseJwt(jwt)
	if err != nil {
		return false
	}

	ok := slices.Contains(roles, claims.Role)
	if !ok {
		return false
	}
	return true
}

func HandleError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}
		var status int
		switch e := err.(type) {
		case *echo.HTTPError:
			return echo.NewHTTPError(e.Code, e.Message)
		case errors.ObjectNotFound, errors.ObjectHasNotSubObjects, errors.NoInProgressReception:
			status = http.StatusNotFound
		case errors.PropertyMissing, errors.PropertyTooSmall, errors.PropertyTooBig, errors.ObjectAlreadyExists,
			errors.WrongPropertyValue, errors.MalformedBody, errors.ReceptionIsNotClosed, errors.ReceptionIsNotInProgress,
			errors.BadPropertyValue, errors.BadParamValue, errors.ParamMissing, errors.StartDateAfterEndDate:
			status = http.StatusBadRequest
		case errors.AccessForbidden:
			status = http.StatusForbidden
		case errors.InvalidCredentials:
			status = http.StatusUnauthorized
		case errors.InternalError:
			status = http.StatusInternalServerError
		default:
			status = http.StatusInternalServerError
		}

		logger.Log.Error(fmt.Sprintf("request ended with error=%s, status code=%v", err, status))
		return c.JSON(status, err)
	}
}
