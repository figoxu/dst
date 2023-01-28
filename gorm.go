package dst

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func ConnectToPG(conn string, opts ...Option) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		o.apply(db)
	}
	return db, nil
}

// NotFound 是否没有找到
func NotFound(err error) bool {
	if err == nil {
		return false
	}

	l := []error{gorm.ErrRecordNotFound, sql.ErrNoRows}

	for _, v := range l {
		if errors.Is(err, v) {
			return true
		}
	}

	// 有时error会被rpc远程传递，变成rpc error，这时只能用字符串判断了
	strList := []string{gorm.ErrRecordNotFound.Error(), sql.ErrNoRows.Error()}
	for _, v := range strList {
		if err.Error() == v {
			return true
		}
	}

	return false
}

func NewGormLogger() gormlogger.Interface {
	return gormlogger.New(log.StandardLogger(), gormlogger.Config{
		SlowThreshold:             time.Second,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  gormlogger.Warn,
	})
}

// Exist 是否存在
func Exist(db *gorm.DB) (bool, error) {
	var n int
	err := db.Select(`1`).Limit(1).Row().Scan(&n)
	if NotFound(err) {
		return false, nil
	}
	return true, err
}
