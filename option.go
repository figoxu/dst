package dst

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Option 是一个连接选项 interface.
type Option interface {
	apply(db *gorm.DB)
}

type optionFunc func(db *gorm.DB)

func (f optionFunc) apply(db *gorm.DB) {
	f(db)
}

// WithConnMaxLifetime sets the maximum amount of time a connection may be reused.
func WithConnMaxLifetime(d time.Duration) Option {
	return optionFunc(func(db *gorm.DB) {
		sqlDb, err := db.DB()
		chk(err)
		sqlDb.SetConnMaxLifetime(d)
	})
}

// WithMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func WithMaxIdleConns(n int) Option {
	return optionFunc(func(db *gorm.DB) {
		sqlDb, err := db.DB()
		chk(err)
		sqlDb.SetMaxIdleConns(n)
	})
}

// WithMaxOpenConns sets the maximum number of open connections to the database.
func WithMaxOpenConns(n int) Option {
	return optionFunc(func(db *gorm.DB) {
		sqlDb, err := db.DB()
		chk(err)
		sqlDb.SetMaxOpenConns(n)
	})
}

// WithPing send db ping repeatedly, avoid db connection being released
func WithPing(d time.Duration) Option {
	fn := func(db *gorm.DB) {
		for {
			time.Sleep(d)
			ctx := context.Background()
			if sqlDb, err := db.DB(); err != nil {
				db.Logger.Warn(ctx, "get sql db err", err)
			} else {
				if err := sqlDb.Ping(); err != nil {
					db.Logger.Warn(ctx, "db ping err", err)
				}
			}
		}
	}
	return optionFunc(func(db *gorm.DB) {
		go fn(db)
	})
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
