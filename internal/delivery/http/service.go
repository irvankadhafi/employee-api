package http

import (
	"github.com/irvankadhafi/employee-api/internal/model"
	"github.com/labstack/echo/v4"
)

// service http service
type service struct {
	employeeUsecase model.EmployeeUsecase
}

// RouteService ..
func RouteService(
	group *echo.Group,
	employeeUsecase model.EmployeeUsecase,
) {
	svc := &service{
		employeeUsecase: employeeUsecase,
	}

	svc.initRoutes(group)
}

func (s *service) initRoutes(group *echo.Group) {
	group.GET("/positions/", s.GetDistinctPositions())

	employeeRoute := group.Group("/employees")
	{
		employeeRoute.POST("/", s.Create())
		employeeRoute.GET("/:employee_id/", s.GetDetail())
		employeeRoute.GET("/", s.SearchEmployees())
		employeeRoute.PUT("/:employee_id/", s.Update())
		employeeRoute.DELETE("/:employee_id/", s.Delete())
	}
}
