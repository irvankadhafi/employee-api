package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/irvankadhafi/employee-api/utils"
	"github.com/labstack/echo/v4"
	"strconv"
	"sync"
)

// validate singleton, it's thread safe and cached the struct validation rules
var validate *validator.Validate

var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		validate = validator.New()
	})
}

type metaInfo struct {
	Size      int  `json:"size"`
	Count     int  `json:"count"`
	CountPage int  `json:"count_page"`
	Page      int  `json:"page"`
	NextPage  int  `json:"next_page"`
	HasMore   bool `json:"has_more"`
}

type paginationResponse[T any] struct {
	Items    []T `json:"items"`
	metaInfo `json:"meta_info"`
}

func toResourcePaginationResponse[T any](page, size int, count int64, items []T) paginationResponse[T] {
	pagination := paginationResponse[T]{
		Items: items,
	}

	if size > 0 {
		pagination.metaInfo = metaInfo{
			Size:      size,
			Count:     int(count),
			CountPage: utils.CalculatePages(int(count), size),
			Page:      page,
			NextPage:  0,
		}
	}

	pagination.metaInfo.HasMore = int(count)-(page*size) > 0
	if pagination.metaInfo.HasMore {
		pagination.metaInfo.NextPage = page + 1
	}

	return pagination
}

type successResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func setSuccessResponse(data any) successResponse {
	return successResponse{
		Success: true,
		Data:    data,
	}
}

type errorResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
}

func setErrorMessage(msg any) errorResponse {
	return errorResponse{
		Success: false,
		Message: msg,
	}
}

// parseQueryParam is a helper function to parse and return an int from a query param or fallback to default.
func parseQueryParam(c echo.Context, param string, defaultValue int) (int, error) {
	paramStr := c.QueryParam(param)
	if paramStr == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(paramStr)
}
