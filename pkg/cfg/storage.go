package cfg

import (
	"time"

	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

type Storage struct {
	PgDsn *PgDsn

	// defines the maximum number of DB connections that can be opened
	// If not set the default value is UNLIMITED.
	MaxOpenConn int

	// defines the maximum number of idle connection that can be present
	// in the pool at a given time. The default value is 2 if not set
	MaxIdleConn int

	// defines the ttl for the DB connection.
	MaxConnTtl time.Duration

	GormV2Conf *GormV2Conf
}

// PgDsn contains standard set of parameters needed to construct a postgres compatible DB connection DSN like AWS RDS
// Refer- https://www.postgresql.org/docs/current/libpq-connect.html
type PgDsn struct {
	Host          string
	Port          int
	Username      string
	Password      string
	Name          string
	SSLMode       string
	SSLCertPath   string
	SSLRootCert   string
	SSLClientCert string
	SSLClientKey  string
	// application name to be set in DB session connection
	AppName string
}

type GormV2Conf struct {
	// LogLevelGormV2 will correspond to the levels defined in this doc https://pkg.go.dev/gorm.io/gorm/logger#LogLevel
	LogLevelGormV2 GormLogLevel
	// Queries with execution time more than SlowQueryLogThreshold will be logged as Slow query.
	SlowQueryLogThreshold time.Duration
	// Set UseInsecureLog to true only if all data can be logged into insecure logs
	// SQL logs may contain values passed as filters. Set this flag cautiously.
	// By default, the flag will be false and always uses Secure logs for SQL
	UseInsecureLog bool

	// Set DisableImplicitPreparedStmt if you want to disable driver(pgx)
	// preparing statement before executing. This should not be used for application code.
	// This is suitable for one-off scripts or batch sql executor to bypass error
	// `prepared statement had x statements, expected 1 (SQLSTATE 42P14)`
	DisableImplicitPreparedStmt bool
}

// LogLevel corresponds to the log level defined in Zap logger pkg. Refer https://pkg.go.dev/go.uber.org/zap/zapcore#Level
type LogLevel string

const (
	DebugLogLevel  LogLevel = "DEBUG"
	InfoLogLevel   LogLevel = "INFO"
	WarnLogLevel   LogLevel = "WARN"
	ErrorLogLevel  LogLevel = "ERROR"
	DPanicLogLevel LogLevel = "DPANIC"
	PanicLogLevel  LogLevel = "PANIC"
	FatalLogLevel  LogLevel = "FATAL"
)

// GormLogLevel corresponds to the log level defined in Zap logger pkg. Refer https://pkg.go.dev/go.uber.org/zap/zapcore#Level
type GormLogLevel string

const (
	SilentGormLogLevel GormLogLevel = "SILENT"
	ErrorGormLogLevel  GormLogLevel = "ERROR"
	WarnGormLogLevel   GormLogLevel = "WARN"
	InfoGormLogLevel   GormLogLevel = "INFO"
)

func GetGORMLogLevel(logLevel GormLogLevel) gormlogger.LogLevel {
	switch logLevel {
	case ErrorGormLogLevel:
		return gormlogger.Error
	case WarnGormLogLevel:
		return gormlogger.Warn
	case InfoGormLogLevel:
		return gormlogger.Info
	case SilentGormLogLevel:
		fallthrough
	default:
		return gormlogger.Silent
	}
}

func GetLogLevel(logLevel LogLevel) zapcore.Level {
	switch logLevel {
	case DebugLogLevel:
		return zapcore.DebugLevel
	case InfoLogLevel:
		return zapcore.InfoLevel
	case WarnLogLevel:
		return zapcore.WarnLevel
	case ErrorLogLevel:
		return zapcore.ErrorLevel
	case DPanicLogLevel:
		return zapcore.DPanicLevel
	case PanicLogLevel:
		return zapcore.PanicLevel
	case FatalLogLevel:
		fallthrough
	default:
		return zapcore.FatalLevel
	}
}
