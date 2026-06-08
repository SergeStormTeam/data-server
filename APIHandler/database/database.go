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
			id UUID PRIMARY KEY,			
			timestamp TIMESTAMPTZ NOT NULL,
			session_id UUID NOT NULL,

			sequence INTEGER NOT NULL,		

            temperature DOUBLE PRECISION,
			humidity DOUBLE PRECISION,
			pressure DOUBLE PRECISION,
			voc DOUBLE PRECISION,
			wind_speed DOUBLE PRECISION,
			co2 DOUBLE PRECISION,
			precipitation DOUBLE PRECISION,

			UNIQUE(session_id, sequence)
        );
    `)
	if err != nil {
		return err
	}

	_, err = newDB.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY,			
			timestamp TIMESTAMPTZ NOT NULL,
			session_id UUID NOT NULL,			

 			message TEXT NOT NULL,
            severity INTEGER NOT NULL
        );
    `)
	if err != nil {
		return err
	}

	db = newDB

	return nil
}
