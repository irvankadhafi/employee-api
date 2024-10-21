package usecase

import (
	"context"
	"github.com/irvankadhafi/employee-api/internal/model"
	"github.com/irvankadhafi/employee-api/utils"
	"github.com/sirupsen/logrus"
	"sync"
)

type employeeUsecase struct {
	employeeRepository model.EmployeeRepository
}

func NewEmployeeUsecase(repository model.EmployeeRepository) model.EmployeeUsecase {
	return &employeeUsecase{employeeRepository: repository}
}

func (e *employeeUsecase) Create(ctx context.Context, input model.CreateEmployeeRequest) (employee *model.Employee, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		logger.Error(err)
		return nil, err
	}

	employee = &model.Employee{
		Name:     input.Name,
		Position: input.Position,
		Salary:   input.Salary,
	}

	if err := e.employeeRepository.Create(ctx, employee); err != nil {
		logger.Error(err)
		return nil, err
	}

	return e.FindByID(ctx, employee.ID)
}

func (e *employeeUsecase) FindByID(ctx context.Context, id int64) (employee *model.Employee, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	employee, err = e.employeeRepository.FindByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if employee == nil {
		return nil, ErrNotFound
	}

	return employee, nil
}

func (e *employeeUsecase) Update(ctx context.Context, employeeID int64, input model.UpdateEmployeeRequest) (employee *model.Employee, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.DumpIncomingContext(ctx),
		"employeeID": employeeID,
		"input":      utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		logger.Error(err)
		return nil, err
	}

	employee, err = e.employeeRepository.FindByID(ctx, employeeID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	employee.Name = input.Name
	employee.Position = input.Position

	if err := e.employeeRepository.Update(ctx, employee); err != nil {
		logger.Error(err)
		return nil, err
	}

	return e.FindByID(ctx, employeeID)
}

func (e *employeeUsecase) DeleteByID(ctx context.Context, employeeID int64) (err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.DumpIncomingContext(ctx),
		"employeeID": employeeID,
	})

	employee, err := e.FindByID(ctx, employeeID)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err := e.employeeRepository.Delete(ctx, employee.ID); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *employeeUsecase) SearchByCriteria(ctx context.Context, searchCriteria model.EmployeeSearchCriteria) (employees []*model.Employee, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":            utils.DumpIncomingContext(ctx),
		"searchCriteria": utils.Dump(searchCriteria),
	})

	ids, count, err := e.searchByPage(ctx, searchCriteria)
	if err != nil {
		logger.Error(err)
		return
	}

	employees = e.findAllByIDs(ctx, ids)
	if len(employees) <= 0 {
		logger.Error(ErrNotFound)
		return
	}

	return employees, count, nil
}

func (e *employeeUsecase) GetDistinctPositions(ctx context.Context) ([]string, error) {
	positions, err := e.employeeRepository.GetDistinctPositions(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if len(positions) <= 0 {
		err = ErrNotFound
		logrus.Error(err)
		return nil, err
	}

	return positions, nil
}

func (e *employeeUsecase) searchByPage(ctx context.Context, searchCriteria model.EmployeeSearchCriteria) (ids []int64, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":            utils.DumpIncomingContext(ctx),
		"searchCriteria": utils.Dump(searchCriteria),
	})

	searchCriteria.SetDefaultValue()
	ids, count, err = e.employeeRepository.SearchByPage(ctx, searchCriteria)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	if len(ids) == 0 || count == 0 {
		return nil, 0, nil
	}

	return ids, count, nil
}

func (e *employeeUsecase) findAllByIDs(ctx context.Context, ids []int64) (employees []*model.Employee) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"ids": ids,
	})

	var wg sync.WaitGroup
	c := make(chan *model.Employee, len(ids))
	for _, id := range ids {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()

			employee, err := e.FindByID(ctx, id)
			if err != nil {
				logger.Error(err)
				return
			}
			c <- employee
		}(id)
	}
	wg.Wait()
	close(c)

	if len(c) <= 0 {
		return
	}

	rs := map[int64]*model.Employee{}
	for employee := range c {
		if employee != nil {
			rs[employee.ID] = employee
		}
	}

	for _, id := range ids {
		if employee, ok := rs[id]; ok {
			employees = append(employees, employee)
		}
	}

	return
}
