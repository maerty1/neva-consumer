package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	dsnEnvName                   = "POSTGRESQL_URL"
	maxConnectionsEnvName        = "POSTGRES_MAX_CONNECTIONS"
	maxConnectionIdleTimeEnvName = "CONNECTION_IDLE_TIME_SEC"

	rootPrefix = "ROOT_"
	zuluPrefix = "ZULU_"

	defaultMaxConnections        = 20
	defaultMaxConnectionIdleTime = 5 * time.Second
)

type PGConfig interface {
	DSN() string
	MaxConnections() int32
	MaxConnectionIdleTime() time.Duration
}

type pgConfig struct {
	dsn                   string
	maxConnections        int32
	maxConnectionIdleTime time.Duration
}

func GetPostgresConfig(dbname string) (PGConfig, error) {
	var prefix string
	switch dbname {
	case "zulu":
		prefix = zuluPrefix
	case "root":
		prefix = rootPrefix
	default:
		return nil, errors.New("неверно указано название бд")

	}
	dsn := os.Getenv(prefix + dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("dsn базы данной не найден")
	}

	maxConnectionsStr := os.Getenv(prefix + maxConnectionsEnvName)
	maxConnections, err := strconv.ParseInt(maxConnectionsStr, 10, 64)
	if err != nil || maxConnections == 0 {
		maxConnections = defaultMaxConnections
	}

	maxConnectionIdleTimeStr := os.Getenv(prefix + maxConnectionIdleTimeEnvName)
	maxConnectionIdleTimeNum, err := strconv.Atoi(maxConnectionIdleTimeStr)
	maxConnectionIdleTime := time.Duration(maxConnectionIdleTimeNum) * time.Second
	if err != nil || maxConnectionIdleTime == 0 {
		maxConnectionIdleTime = defaultMaxConnectionIdleTime
	}

	return &pgConfig{
		dsn:                   dsn,
		maxConnections:        int32(maxConnections),
		maxConnectionIdleTime: maxConnectionIdleTime,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}

func (cfg *pgConfig) MaxConnections() int32 {
	return cfg.maxConnections
}

func (cfg *pgConfig) MaxConnectionIdleTime() time.Duration {
	return cfg.maxConnectionIdleTime
}

func NewPGConfig(dsn string, maxConnections int64, maxConnectionIdleTime int64) PGConfig {
	return &pgConfig{
		dsn:                   dsn,
		maxConnections:        int32(maxConnections),
		maxConnectionIdleTime: time.Duration(maxConnectionIdleTime) * time.Second,
	}
}
