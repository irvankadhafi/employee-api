package http

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	ErrInvalidArgument      = echo.NewHTTPError(http.StatusBadRequest, setErrorMessage("invalid argument"))
	ErrInternal             = echo.NewHTTPError(http.StatusInternalServerError, setErrorMessage("internal system error"))
	ErrNotFound             = echo.NewHTTPError(http.StatusNotFound, setErrorMessage("record not found"))
	ErrEmployeeAlreadyExist = echo.NewHTTPError(http.StatusBadRequest, setErrorMessage("employee already exist"))
)

// httpValidationOrInternalErr return valdiation or internal error
func httpValidationOrInternalErr(err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// Jika tidak ada kesalahan validasi, mengembalikan kesalahan internal
		return ErrInternal
	}

	fields := make(map[string]string)
	for _, validationError := range validationErrors {
		tag := validationError.Tag()
		fields[validationError.Field()] = fmt.Sprintf("Failed on the '%s' tag", tag)
	}

	return echo.NewHTTPError(http.StatusBadRequest, setErrorMessage(fields))
}
