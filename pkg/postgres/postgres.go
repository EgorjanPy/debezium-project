package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host               string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost"`
	Port               uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432"`
	User               string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-default:"postgres"`
	Password           string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASS" env-default:"postgres"`
	DataBase           string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-default:"postgres"`
	MaxOpenConnections int32  `yaml:"POSTGRES_MAX_OPEN_CONNECTIONS" env:"POSTGRES_MAX_OPEN_CONNECTIONS" env-default:"10"`
	MinOpenConnections int32  `yaml:"POSTGRES_MIN_OPEN_CONNECTIONS" env:"POSTGRES_MIN_OPEN_CONNECTIONS" env-default:"5"`
}
type DataBase struct {
	Pool *pgxpool.Pool
}

// postgres://username:password@localhost:5432/database_name
func New(ctx context.Context, config Config) (*DataBase, error) {
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DataBase,
	))
	if err != nil {
		return nil, fmt.Errorf("postgres.New parse config: %w", err)
	}

	// Настраиваем параметры пула
	poolConfig.MaxConns = config.MaxOpenConnections
	poolConfig.MinConns = config.MinOpenConnections

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres.New ping failed: %w", err)
	}

	return &DataBase{Pool: pool}, nil
}
func (db *DataBase) Close() {
	db.Pool.Close()
}
