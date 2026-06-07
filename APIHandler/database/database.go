package database

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newDB, err := pgxpool.New(
		ctx,
		os.Getenv("POSTGRESQL_ADDRESS"),
	)

	if err != nil {
		return err
	}

	err = newDB.Ping(ctx)
	if err != nil {
		return err
	}

	_, err = newDB.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS data (
			session_id UUID NOT NULL,
			timestamp TIMESTAMPZ NOT NULL,
			sequence INTEGER NOT NULL,		

            timestamp TIMESTAMPTZ NOT NULL,
            temperature DOUBLE PRECISION,
			humidity DOUBLE PRECISION,
			pressure DOUBLE PRECISION,
			voc DOUBLE PRECISION,
			wind_speed DOUBLE PRECISION,
			co2 DOUBLE PRECISION,
			precipitation DOUBLE PRECISION,

			PRIMARY KEY (session_id, sequence)
        );
    `)

	db = newDB

	return nil
}
