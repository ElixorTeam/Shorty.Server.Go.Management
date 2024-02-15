package utils

import (
	"fmt"
	"time"

	"Shorty.Server.Go.Mangment/internal/config"
	"Shorty.Server.Go.Mangment/internal/datasources/drivers"
	"github.com/jmoiron/sqlx"
)

func SetupPostgresConnection() (*sqlx.DB, error) {

	databaseConfig := drivers.SQLXConfig{
		DriverName:     config.AppConfig.DBDriver,
		DataSourceName: fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", config.AppConfig.DBHost, config.AppConfig.DBHost, config.AppConfig.DBUser, config.AppConfig.DBName, config.AppConfig.DBPassword),
		MaxOpenConns:   100,
		MaxIdleConns:   10,
		MaxLifetime:    15 * time.Minute,
	}

	conn, err := databaseConfig.InitializeSQLXDatabase()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
