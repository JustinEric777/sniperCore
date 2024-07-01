package db

import (
	"errors"
	"fmt"
	"github.com/sniperCore/core/log"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"
)

const (
	defaultTimeout = 10
	defaultCharset = "utf8mb4"
)

func NewDB(config *DbConfig) (*gorm.DB, error) {
	//获取数据库连接
	db, err := newConn(config.Driver, config.BaseConfig, config.OptionConfig)
	if err != nil {
		return nil, err
	}
	//设置数据库相关信息
	err = setDBOption(db, config.OptionConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newConn(driver string, base DbBaseConfig, option DbOptionConfig) (*gorm.DB, error) {
	dsn := formatDSN(driver, base, option)
	if dsn == "" {
		return nil, errors.New(fmt.Sprintf("missing db driver %s or db config", driver))
	}

	newLogger := New(log.GetLogger(log.SingletonMain), DefaultGormLoggerConfig)
	switch driver {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{
			QueryFields: true,
			Logger:      newLogger,
		})
	case "postgres":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite3":
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "mssql":
		return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	case "clickhouse":
		return gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	}

	return nil, errors.New("not found db driver")
}

func setDBOption(db *gorm.DB, option DbOptionConfig) error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(option.MaxIdle)
	sqlDb.SetMaxIdleConns(option.MaxConns)
	sqlDb.SetConnMaxIdleTime(time.Duration(option.IdleTimeout))
	sqlDb.SetConnMaxLifetime(time.Duration(option.ConnectTimeout))

	return nil
}

func formatDSN(driver string, base DbBaseConfig, option DbOptionConfig) string {
	switch driver {
	case "mysql":
		return formatMysqlDSN(base, option)
	case "postgres":
		return formatPostgresDSN(base, option)
	case "sqlite3":
		return formatSqlite3DSN(base, option)
	case "mssql":
		return formatMssqlDSN(base, option)
	case "clickhouse":
		return formatClickHouseDSN(base, option)
	}
	return ""
}

func formatMysqlDSN(base DbBaseConfig, option DbOptionConfig) string {
	port := getPortOrDefault(base.Port, 3306)
	charset := option.Charset
	if charset == "" {
		charset = defaultCharset
	}
	timeout := option.ConnectTimeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&charset=%s&parseTime=true&loc=Local",
		base.Username, base.Password, base.Host, port, base.DBName, timeout, charset)
}

func formatPostgresDSN(base DbBaseConfig, option DbOptionConfig) string {
	port := getPortOrDefault(base.Port, 5432)
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		base.Host, port, base.Username, base.DBName, base.Password)
}

func formatSqlite3DSN(base DbBaseConfig, option DbOptionConfig) string {
	return base.DBName
}

func formatMssqlDSN(base DbBaseConfig, option DbOptionConfig) string {
	port := getPortOrDefault(base.Port, 1433)
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		base.Username, base.Password, base.Host, port, base.DBName)
}

func formatClickHouseDSN(base DbBaseConfig, option DbOptionConfig) string {
	port := getPortOrDefault(base.Port, 9000)
	return fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s&read_timeout=10&write_timeout=20",
		base.Host, port, base.DBName, base.Username, base.Password)
}

func getPortOrDefault(port int, defaultPort int) int {
	if port == 0 {
		return defaultPort
	}
	return port
}
