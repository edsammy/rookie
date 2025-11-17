package db

import (
	"context"
	"database/sql"
	"time"
)

type BloodGlucoseReading struct {
	UserID     string
	SampleTime time.Time
	MgDL       sql.NullInt64
}

const upsertBloodGlucose = `
	INSERT INTO blood_glucose (user_id, sample_time, mg_dl)
	VALUES (?, ?, ?)
	ON CONFLICT(user_id, sample_time) DO UPDATE SET
		mg_dl = excluded.mg_dl;
	`

func UpsertBloodGlucose(ctx context.Context, reading BloodGlucoseReading) error {
	_, err := DB.ExecContext(ctx, upsertBloodGlucose, reading.UserID, reading.SampleTime, reading.MgDL)
	return err
}

func UpsertBloodGlucoseBatch(ctx context.Context, readings []BloodGlucoseReading) error {
	if len(readings) == 0 {
		return nil
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, upsertBloodGlucose)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, r := range readings {
		if _, err := stmt.ExecContext(ctx, r.UserID, r.SampleTime, r.MgDL); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
