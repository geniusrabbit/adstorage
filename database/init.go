package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/demdxx/xtypes"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/context/ctxdatabase"
)

type openFnk func(dsn string) gorm.Dialector

var dialectors = map[string]openFnk{}

type DB = gorm.DB

// Connect to database
func Connect(ctx context.Context, connection string, debug bool) (*DB, error) {
	var (
		i      = strings.Index(connection, "://")
		driver = connection[:i]
	)
	if driver == "mysql" {
		connection = connection[i+3:]
	}
	openDriver := dialectors[driver]
	if openDriver == nil {
		return nil, fmt.Errorf(`unsupported database driver %s`, driver)
	}
	db, err := gorm.Open(openDriver(connection), &gorm.Config{SkipDefaultTransaction: true})
	if err == nil && debug {
		db = db.Debug()
	}
	return db, err
}

// WithDatabase puts databases to context
func WithDatabase(ctx context.Context, master, slave *gorm.DB) context.Context {
	return ctxdatabase.WithDatabase(ctx, master, slave)
}

// ListOfDialects returns list of available DB drivers
func ListOfDialects() []string {
	return xtypes.Map[string, openFnk](dialectors).Keys()
}
