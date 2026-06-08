package database

import (
	"context"
	"time"

	"github.com/WeatherGod3218/serge-api-handler/logging"

	"github.com/jackc/pgx/v5"

	"github.com/sirupsen/logrus"
)

type DBDataEntry struct {
	RecordId      string   `json:"id"`
	Timestamp     float64  `json:"timestamp"`
	SessionId     string   `json:"session_id"`
	Sequence      int      `json:"sequence"`
	CO2           *float64 `json:"co2"`
	Humidity      *float64 `json:"humidity"`
	Precipitation *float64 `json:"precipitation"`
	Pressure      *float64 `json:"pressure"`
	VOC           *float64 `json:"voc"`
	WindSpeed     *float64 `json:"wind_speed"`
}

type DBEventEntry struct {
	RecordId  string  `json:"id"`
	Timestamp float64 `json:"timestamp"`
	SessionId string  `json:"session_id"`
	Message   string  `json:"message"`
	Severity  int     `json:"severity"`
}

type DatabaseResponse struct {
	RecordId string `json:"id"`
}

type DatabaseBackup struct {
	Events []DBEventEntry `json:"events"`
	Data   []DBDataEntry  `json:"data"`
}

// func mustParseUUID(s string) pgtype.UUID {
// 	var u pgtype.UUID
// 	err := u.Scan(s)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return u
// }

func AddEventsToDatabase(new_rows []DBEventEntry) ([]DatabaseResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows := make([][]any, len(new_rows))
	returned_rows := make([]DatabaseResponse, len(new_rows))

	for i, entry := range new_rows {
		rows[i] = []any{
			entry.RecordId,
			time.Unix(int64(entry.Timestamp), 0).UTC(),
			entry.SessionId,
			entry.Message,
			entry.Severity,
		}

		returned_rows[i] = DatabaseResponse{
			entry.RecordId,
		}
	}

	transaction, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = transaction.Rollback(ctx)
	}()

	transaction.Exec(ctx, `CREATE TEMP TABLE events_staging (LIKE events INCLUDING ALL) ON COMMIT DROP`)

	_, err = transaction.CopyFrom(
		ctx,
		pgx.Identifier{"events_staging"},
		[]string{"id", "timestamp", "session_id", "message", "severity"},
		pgx.CopyFromRows(rows),
	)

	transaction.Exec(ctx, `
		INSERT INTO events
		SELECT * FROM events_staging
		ON CONFLICT DO NOTHING
	`)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err}).Warn("Failed to copy events into the database!")
		return nil, err
	}

	err = transaction.Commit(ctx)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err}).Warn("Failed to commit events to the database!")
		return nil, err
	}

	return returned_rows, nil
}

func AddDataToDatabase(new_rows []DBDataEntry) ([]DatabaseResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows := make([][]any, len(new_rows))
	returned_rows := make([]DatabaseResponse, len(new_rows))

	for i, entry := range new_rows {
		rows[i] = []any{
			entry.RecordId,
			time.Unix(int64(entry.Timestamp), 0).UTC(),
			entry.SessionId,
			entry.Sequence,
			entry.CO2,
			entry.Humidity,
			entry.Precipitation,
			entry.Pressure,
			entry.VOC,
			entry.WindSpeed,
		}

		returned_rows[i] = DatabaseResponse{
			entry.RecordId,
		}
	}

	transaction, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = transaction.Rollback(ctx)
	}()

	transaction.Exec(ctx, `CREATE TEMP TABLE data_staging (LIKE data INCLUDING ALL) ON COMMIT DROP`)

	_, err = transaction.CopyFrom(
		ctx,
		pgx.Identifier{"data_staging"},
		[]string{"id", "timestamp", "session_id", "sequence", "co2", "humidity", "precipitation", "pressure", "voc", "wind_speed"},
		pgx.CopyFromRows(rows),
	)

	transaction.Exec(ctx, `
		INSERT INTO data
		SELECT * FROM data_staging
		ON CONFLICT DO NOTHING
	`)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err}).Warn("Failed to copy data into the database!")
		return nil, err
	}

	err = transaction.Commit(ctx)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err}).Warn("Failed to commit data to the database!")
		return nil, err
	}

	return returned_rows, nil
}
