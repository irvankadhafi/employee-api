package http

import (
	"github.com/irvankadhafi/employee-api/internal/model"
	"github.com/irvankadhafi/employee-api/internal/usecase"
	"github.com/irvankadhafi/employee-api/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *service) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		req := model.CreateEmployeeRequest{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}

		createdProduct, err := s.employeeUsecase.Create(ctx, req)
		switch err {
		case nil:
			break
		case usecase.ErrDuplicateEmployee:
			return ErrEmployeeAlreadyExist
		default:
			logrus.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusCreated, setSuccessResponse(createdProduct))
	}
}

func (s *service) GetDetail() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		productID := utils.StringToInt64(c.Param("employee_id"))
		logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"employee_id": productID,
		})

		product, err := s.employeeUsecase.FindByID(ctx, productID)
		if err != nil {
			logger.Error(err)
			return ErrInternal
		}

		switch err {
		case nil:
			break
		case usecase.ErrNotFound:
			return ErrNotFound
		case usecase.ErrDuplicateEmployee:
			return ErrEmployeeAlreadyExist
		default:
			logrus.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusOK, setSuccessResponse(product))
	}
}

func (s *service) GetList() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// Parse query parameters with default values
		page, err := parseQueryParam(c, "page", 1)
		if err != nil {
			logrus.WithError(err).Error("failed to parse page")
			return ErrInvalidArgument
		}

		limit, err := parseQueryParam(c, "limit", 10)
		if err != nil {
			logrus.WithError(err).Error("failed to parse limit")
			return ErrInvalidArgument
		}

		query := c.QueryParam("query")
		dir := c.QueryParam("dir")
		sort := c.QueryParam("sort")

		// Define search criteria
		searchCriteria := model.EmployeeSearchCriteria{
			Query:   query,
			Page:    int64(page),
			Size:    int64(limit),
			SortBy:  sort,
			SortDir: dir,
		}

		employees, count, err := s.employeeUsecase.SearchByCriteria(ctx, searchCriteria)
		if err != nil {
			logrus.WithError(err).Error("failed to retrieve employees")
			return c.JSON(http.StatusBadRequest, "Error retrieving employees")
		}

		logrus.WithFields(logrus.Fields{
			"page":  page,
			"limit": limit,
		}).Info("success retrieving employees")

		return c.JSON(http.StatusOK, toResourcePaginationResponse(page, limit, count, employees))
	}
}

func (s *service) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		req := model.UpdateEmployeeRequest{}
		if err := c.Bind(&req); err != nil {
			logrus.Error(err)
			return ErrInvalidArgument
		}
		employeeID := utils.StringToInt64(c.Param("employee_id"))

		newEmployee, err := s.employeeUsecase.Update(ctx, employeeID, req)
		switch err {
		case nil:
			break
		case usecase.ErrNotFound:
			return ErrNotFound
		default:
			logrus.Error(err)
			return ErrInternal
		}

		return c.JSON(http.StatusCreated, setSuccessResponse(newEmployee))
	}
}

func (s *service) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		employeeID := utils.StringToInt64(c.Param("employee_id"))

		if err := s.employeeUsecase.DeleteByID(ctx, employeeID); err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"employeeID": employeeID,
			}).Error(err)

			return ErrInternal
		}

		return c.JSON(http.StatusOK, setSuccessResponse(employeeID))
	}
}
