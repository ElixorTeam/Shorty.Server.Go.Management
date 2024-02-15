package drivers

import (
	"fmt"
	"time"

	"Shorty.Server.Go.Management/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type SQLXConfig struct {
	DriverName     string
	DataSourceName string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetime    time.Duration
}

func (config *SQLXConfig) InitializeSQLXDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	logger.Info(fmt.Sprintf("Setting maximum number of open connections to %d", config.MaxOpenConns), nil)
	db.SetMaxOpenConns(config.MaxOpenConns)

	logger.Info(fmt.Sprintf("Setting maximum number of idle connections to %d", config.MaxIdleConns), nil)
	db.SetMaxIdleConns(config.MaxIdleConns)

	logger.Info(fmt.Sprintf("Setting maximum lifetime for a connection to %s", config.MaxLifetime), nil)
	db.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}
