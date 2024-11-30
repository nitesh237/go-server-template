package storage

import (
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/nitesh237/go-server-template/pkg/cfg"
	logpkg "github.com/nitesh237/go-server-template/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

const (
	PostgresSQLSchema = "postgresql"
	PostgresDriver    = "postgres"
)

const (
	DBConnSSLMode     = "sslmode"
	DBConnSSLRootCert = "sslrootcert"
	DBConnSSLKey      = "sslkey"
	DBConnSSLCert     = "sslcert"
	DBConnAppName     = "application_name"
	StatementTimeout  = "statement_timeout"
)

const (
	DBSSLModeDisable    = "disable"
	DBSSLModeVerifyFull = "verify-full"
	DBSSLModeRequired   = "required"
)

func NewPostgresDB(dbConf *cfg.Storage, loger logpkg.Logger) (*gorm.DB, error) {
	connString, err := GetPgDbDsnStringFromDBConf(dbConf)
	if err != nil {
		return nil, err
	}

	gormConfig := &gorm.Config{
		NowFunc: pgnow,
		Logger: gormlogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gormlogger.Config{
			SlowThreshold: dbConf.GormV2Conf.SlowQueryLogThreshold,
			LogLevel:      cfg.GetGORMLogLevel(dbConf.GormV2Conf.LogLevelGormV2),
		}),
	}

	if zapLogger, ok := loger.(logpkg.ZapLogger); ok {
		gormConfig.Logger = zapgorm2.New(zapLogger.Unwrap()).LogMode(gormlogger.LogLevel(cfg.GetGORMLogLevel(dbConf.GormV2Conf.LogLevelGormV2)))
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connString,
		PreferSimpleProtocol: dbConf.GormV2Conf.DisableImplicitPreparedStmt,
	}), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if dbConf.MaxOpenConn != 0 {
		sqlDB.SetMaxOpenConns(dbConf.MaxOpenConn)
	}
	if dbConf.MaxIdleConn != 0 {
		sqlDB.SetMaxIdleConns(dbConf.MaxIdleConn)
	}
	if dbConf.MaxConnTtl != 0 {
		sqlDB.SetConnMaxLifetime(dbConf.MaxConnTtl)
	}

	return db, nil
}

func GetPgDbDsnStringFromDBConf(pgdbConf *cfg.Storage) (string, error) {
	connStr, err := GetPgDbDsnString(pgdbConf.PgDsn)
	if err != nil {
		return "", err
	}

	return connStr, nil
}

func GetPgDbDsnString(dsn *cfg.PgDsn) (string, error) {
	uri := GetPgDsnUrl(dsn)

	connStr, err := url.QueryUnescape(uri.String())
	if err != nil {
		return "", err
	}

	return connStr, nil
}

func GetPgDsnUrl(pgDsnConf *cfg.PgDsn) *url.URL {
	queryVals := url.Values{}
	queryVals.Set(DBConnSSLMode, pgDsnConf.SSLMode)

	if pgDsnConf.AppName != "" {
		queryVals.Set(DBConnAppName, pgDsnConf.AppName)
	}

	host := os.Getenv("DB_HOST")
	if host != "" {
		pgDsnConf.Host = host
	}

	uri := &url.URL{
		Scheme:   PostgresSQLSchema,
		User:     url.UserPassword(pgDsnConf.Username, pgDsnConf.Password),
		Host:     net.JoinHostPort(pgDsnConf.Host, strconv.Itoa(pgDsnConf.Port)),
		Path:     pgDsnConf.Name,
		RawQuery: queryVals.Encode(),
	}

	return uri
}

func pgnow() time.Time {
	// postgres supports only microsecond precision.
	// Hence, we round off now() value from nanoseconds to microseconds precision.
	// https://github.com/go-gorm/gorm/issues/3232
	return time.Unix(0, time.Now().UnixNano()/int64(time.Microsecond)*1000)
}
