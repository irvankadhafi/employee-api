package repository

import (
	"context"
	"fmt"
	"github.com/irvankadhafi/employee-api/cacher"
	"github.com/irvankadhafi/employee-api/internal/config"
	"github.com/irvankadhafi/employee-api/internal/model"
	"github.com/irvankadhafi/employee-api/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type employeeRepository struct {
	db           *gorm.DB
	cacheManager cacher.CacheManager
}

func NewEmployeeRepository(db *gorm.DB, cacheManager cacher.CacheManager) model.EmployeeRepository {
	return &employeeRepository{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (e *employeeRepository) FindByID(ctx context.Context, id int64) (*model.Employee, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	cacheKey := e.newCacheKeyByID(id)
	if !config.DisableCaching() {
		reply, mu, err := findFromCacheByKey[*model.Employee](e.cacheManager, cacheKey)
		defer cacher.SafeUnlock(mu)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		if mu == nil {
			return reply, nil
		}
	}

	employee := &model.Employee{}
	err := e.db.WithContext(ctx).Take(employee, "id = ?", id).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNil(e.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	err = e.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.Dump(employee)))
	if err != nil {
		logger.Error(err)
	}

	return employee, nil
}

func (e *employeeRepository) Update(ctx context.Context, employee *model.Employee) (err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"employee": utils.Dump(employee),
	})

	err = e.db.WithContext(ctx).Model(&model.Employee{}).
		Where("id = ?", employee.ID).Select("name", "position", "salary").
		Updates(employee).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	if err := e.cacheManager.DeleteByKeys([]string{
		e.newCacheKeyByID(employee.ID),
	}); err != nil {
		logger.Error(err)
	}

	return nil
}

func (e *employeeRepository) Delete(ctx context.Context, employeeID int64) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.DumpIncomingContext(ctx),
		"employeeID": employeeID,
	})

	employee := model.Employee{}

	err := employee.DeletedAt.Scan(time.Now())
	if err != nil {
		return err
	}

	err = e.db.WithContext(ctx).Model(&model.Employee{}).
		Where("id = ?", employeeID).
		Updates(employee).
		Error
	if err != nil {
		logger.Error(err)
		return err
	}

	err = e.cacheManager.DeleteByKeys([]string{
		e.newCacheKeyByID(employeeID),
	})
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (e *employeeRepository) SearchByPage(ctx context.Context, searchCriteria model.EmployeeSearchCriteria) (ids []int64, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":            utils.DumpIncomingContext(ctx),
		"searchCriteria": utils.Dump(searchCriteria),
	})

	count, err = e.countAll(ctx, searchCriteria)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	if count <= 0 {
		return nil, 0, nil
	}

	ids, err = e.findAllIDsByCriteria(ctx, searchCriteria)
	switch err {
	case nil:
		return ids, count, nil
	case gorm.ErrRecordNotFound:
		return nil, 0, nil
	default:
		logger.Error(err)
		return nil, 0, err
	}
}

func (e *employeeRepository) findAllIDsByCriteria(ctx context.Context, criteria model.EmployeeSearchCriteria) ([]int64, error) {
	var scopes []func(*gorm.DB) *gorm.DB
	scopes = append(scopes, scopeByPageAndLimit(criteria.Page, criteria.Size))

	// Add LIKE query for name if provided
	if criteria.Name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name ILIKE ?", "%"+criteria.Name+"%")
		})
	}

	// Add filter for position if provided
	if criteria.Position != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("position = ?", criteria.Position)
		})
	}

	var ids []int64
	err := e.db.WithContext(ctx).
		Model(model.Employee{}).
		Scopes(scopes...).
		Order(fmt.Sprintf("%s %s", criteria.SortBy, criteria.SortDir)).
		Pluck("id", &ids).Error

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      utils.DumpIncomingContext(ctx),
			"criteria": utils.Dump(criteria),
		}).Error(err)
		return nil, err
	}

	return ids, nil
}

func (e *employeeRepository) Create(ctx context.Context, employee *model.Employee) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"employee": utils.Dump(employee),
	})

	err := e.db.WithContext(ctx).Create(employee).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *employeeRepository) GetDistinctPositions(ctx context.Context) ([]string, error) {
	var positions []string
	err := e.db.WithContext(ctx).Debug().
		Model(&model.Employee{}).
		Select("DISTINCT position").
		Pluck("position", &positions).Error

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return nil, err
	}

	return positions, nil
}

func (e *employeeRepository) countAll(ctx context.Context, criteria model.EmployeeSearchCriteria) (int64, error) {
	var scopes []func(*gorm.DB) *gorm.DB

	// Add LIKE query for name if provided
	if criteria.Name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name ILIKE ?", "%"+criteria.Name+"%")
		})
	}

	// Add filter for position if provided
	if criteria.Position != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("position = ?", criteria.Position)
		})
	}

	var count int64
	err := e.db.WithContext(ctx).Model(model.Employee{}).
		Scopes(scopes...).
		Count(&count).
		Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":      utils.DumpIncomingContext(ctx),
			"criteria": utils.Dump(criteria),
		}).Error(err)
		return 0, err
	}

	return count, nil
}

func (e *employeeRepository) newCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:employee:id:%d", id)
}
