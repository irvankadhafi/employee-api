package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type EmployeeUsecase interface {
	Create(ctx context.Context, input CreateEmployeeRequest) (employee *Employee, err error)
	FindByID(ctx context.Context, id int64) (employee *Employee, err error)
	Update(ctx context.Context, employeeID int64, input UpdateEmployeeRequest) (employee *Employee, err error)
	DeleteByID(ctx context.Context, employeeID int64) (err error)
	SearchByCriteria(ctx context.Context, searchCriteria EmployeeSearchCriteria) (employees []*Employee, count int64, err error)
	GetDistinctPositions(ctx context.Context) ([]string, error)
}

type EmployeeRepository interface {
	Create(ctx context.Context, employee *Employee) error
	FindByID(ctx context.Context, id int64) (*Employee, error)
	Update(ctx context.Context, employee *Employee) (err error)
	Delete(ctx context.Context, id int64) error
	SearchByPage(ctx context.Context, searchCriteria EmployeeSearchCriteria) (ids []int64, count int64, err error)
	GetDistinctPositions(ctx context.Context) ([]string, error)
}

// Employee :nodoc:
type Employee struct {
	ID        int64          `json:"id" gorm:"<-:create; primary_key;AUTO_INCREMENT"`
	Name      string         `json:"name"`
	Position  string         `json:"position"`
	Salary    float64        `json:"salary"`
	CreatedAt *time.Time     `json:"created_at" gorm:"->;<-:create"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// CreateEmployeeRequest DTO for creating a new employee
type CreateEmployeeRequest struct {
	Name     string  `json:"name" validate:"required"`
	Position string  `json:"position" validate:"required"`
	Salary   float64 `json:"salary" validate:"required"`
}

func (c *CreateEmployeeRequest) Validate() error {
	return validate.Struct(c)
}

// UpdateEmployeeRequest DTO for updating an employee
type UpdateEmployeeRequest struct {
	Name     string  `json:"name,omitempty"`
	Position string  `json:"position,omitempty"`
	Salary   float64 `json:"salary,omitempty"`
}

func (c *UpdateEmployeeRequest) Validate() error {
	return validate.Struct(c)
}

// EmployeeSearchCriteria :nodoc:
type EmployeeSearchCriteria struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Page     int64  `json:"page"`
	Size     int64  `json:"size"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
}

// SetDefaultValue will set default value for page and size if zero
func (c *EmployeeSearchCriteria) SetDefaultValue() {
	if c.Page == 0 {
		c.Page = 1
	}
	if c.Size == 0 {
		c.Size = 10
	}
	if c.SortBy == "" {
		c.SortBy = "created_at"
	}
	if c.SortDir == "" {
		c.SortDir = "desc"
	}
}
