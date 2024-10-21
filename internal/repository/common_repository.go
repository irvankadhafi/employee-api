package repository

import (
	"encoding/json"
	"github.com/go-redsync/redsync/v4"
	"github.com/irvankadhafi/employee-api/cacher"
	"github.com/irvankadhafi/employee-api/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func storeNil(ck cacher.CacheManager, key string) {
	err := ck.StoreNil(key)
	if err != nil {
		logrus.Error(err)
	}
}

// scopeByPageAndLimit is a helper function to apply pagination on gorm query.
// it takes in 2 input as page and limit and returns a scope function
// that can be passed to gorm's db.Scopes method
// It is reusable to apply pagination on any query where it is needed
func scopeByPageAndLimit(page, limit int64) func(d *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db.Offset(utils.Offset(int(page), int(limit))).Limit(int(limit)) }
}

func findFromCacheByKey[T any](cache cacher.CacheManager, key string) (item T, mutex *redsync.Mutex, err error) {
	var cachedData any

	cachedData, mutex, err = cache.GetOrLock(key)
	if err != nil || cachedData == nil {
		return
	}

	cachedDataByte, _ := cachedData.([]byte)
	if cachedDataByte == nil {
		return
	}

	if err = json.Unmarshal(cachedDataByte, &item); err != nil {
		return
	}

	return
}

func withSize(size int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(size))
	}
}
