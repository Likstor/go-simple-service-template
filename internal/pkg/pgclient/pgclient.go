package pgclient

import (
	"context"
	"fmt"
	"service/internal/pkg/logs"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrFailedConnectionToDatabase = fmt.Errorf("failed connection to database")
)

type Config struct {
	Username string
	Password string
	Host string
	Port string
	Database string
}

const op = "pgclient.NewClient"

func NewClient(ctx context.Context, connAttempts int, config Config) (*pgxpool.Pool, error) {
	logs.Info(
		ctx, 
		"pool creation started",
		op,
	)
	
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		connAttempts = 0
	}

	for i := 0; i < connAttempts; i++ {
		ctx, cancel := context.WithTimeout(ctx, 5 * time.Second)
		defer cancel()
		
		logs.Info(
			ctx,
			fmt.Sprintf("attempt %d to connect to the database", i),
			op,
		)
		err = pool.Ping(ctx)
		if err != nil {
			logs.Warn(
				ctx,
				"failed to connect to the database",
				op,
			)
			time.Sleep(5 * time.Second)
			continue
		}

		logs.Info(
			ctx,
			"successful connection to the database, pool creation successful",
			op,
		)
		
		return pool, nil
	}

	logs.Error(
		ctx,
		"pool creation failed",
		op,
	)
	return nil, ErrFailedConnectionToDatabase
}