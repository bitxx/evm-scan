package runtime

import (
	"gorm.io/gorm"
	"sync"
)

type Application struct {
	dbs map[string]*gorm.DB
	mux sync.RWMutex
}

// NewConfig 默认值
func NewConfig() *Application {
	return &Application{
		dbs: make(map[string]*gorm.DB),
	}
}

// SetDb 设置对应key的db
func (e *Application) SetDb(key string, db *gorm.DB) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.dbs[key] = db
}

// GetDb 获取所有map里的db数据
func (e *Application) GetDb() map[string]*gorm.DB {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.dbs
}

// GetDbByKey 根据key获取db
func (e *Application) GetDbByKey(key string) *gorm.DB {
	e.mux.Lock()
	defer e.mux.Unlock()
	if db, ok := e.dbs["*"]; ok {
		return db
	}
	return e.dbs[key]
}
