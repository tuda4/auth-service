package response

import (
	"github.com/labstack/echo/v4"
	"github.com/tuda4/mb-backend/internal/errors"
)

func ResponseError(c echo.Context, code int, err *errors.BaseError) error {
	return c.JSON(code, err)
}

func ResponseSuccess(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, data)
}

type SuccessResponse[T any] struct {
	Code int `json:"code"`
	Data T   `json:"data"`
}
