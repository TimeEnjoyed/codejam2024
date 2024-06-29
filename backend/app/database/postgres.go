package database

import (
	"codejam.io/config"
	"codejam.io/logging"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var logger = logging.NewLogger(logging.Options{Name: "Database", Level: logging.INFO})

var Pool *pgxpool.Pool // concurrency safe

type Postgres struct {
	Config *config.DBConfig
	Pool   *pgxpool.Pool
}

func Initialize(config config.DBConfig) {
	logger.Info("Connecting to database")

	pool, err := pgxpool.New(context.Background(), config.Url)
	if err != nil {
		logger.Critical("Unable to initialise new pxpool.Pool: %+v\n", err)
		os.Exit(1)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logger.Critical("Unable to connect to database: %+v\n", err)
		os.Exit(1)
	}

	Pool = pool
	logger.Info("Connected to database")
}

func GetRow[T any](query string, args ...any) (T, error) {
	var result T
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		logger.Error("acquire conn error %v", err)
		return result, err
	}
	defer conn.Release()

	// Look at the output value passed in, and get pointers to all the individual fields.  Those will be the Scan params
	row, err := conn.Query(context.Background(),
		query,
		args...)
	if err == nil {
		result, err = pgx.CollectOneRow(row, pgx.RowToStructByName[T])
	} else {
		logger.Error("Error executing query: error %v, query: %v", err, query)
	}
	return result, err
}

func GetRows[T any](query string, args ...any) ([]T, error) {
	var result []T
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		logger.Error("acquire conn error %v", err)
		return result, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(),
		query,
		args...)
	if err == nil {
		result, err = pgx.CollectRows(rows, pgx.RowToStructByName[T])
		if err != nil {
			logger.Error("CollectRows error: %v", err)
		}
	} else {
		logger.Error("Error executing query: error %v, query: %v", err, query)
	}

	return result, nil
}

func Execute(query string, args ...any) error {
	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		logger.Error("acquire conn error %v", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Query(context.Background(),
		query,
		args...)
	if err != nil {
		logger.Error("Query Execute Error: %v", err)
	}

	return nil
}
