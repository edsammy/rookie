package db

import (
	"context"
	"database/sql"
	"time"
)

type HeartRateReading struct {
	UserID     string
	SampleTime time.Time
	BPM        sql.NullInt64
	HrvRMSSD   sql.NullFloat64
	HrvSDNN    sql.NullFloat64
}

const upsertHeartRate = `
	INSERT INTO heart_rate (user_id, sample_time, bpm, hrv_rmssd, hrv_sdnn)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(user_id, sample_time) DO UPDATE SET
		bpm = excluded.bpm,
		hrv_rmssd = excluded.hrv_rmssd,
		hrv_sdnn = excluded.hrv_sdnn;
	`

func UpsertHeartRate(ctx context.Context, reading HeartRateReading) error {
	_, err := DB.ExecContext(ctx, upsertHeartRate,
		reading.UserID,
		reading.SampleTime,
		reading.BPM,
		reading.HrvRMSSD,
		reading.HrvSDNN,
	)
	return err
}

func UpsertHeartRateBatch(ctx context.Context, readings []HeartRateReading) error {
	if len(readings) == 0 {
		return nil
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, upsertHeartRate)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, r := range readings {
		if _, err := stmt.ExecContext(ctx,
			r.UserID,
			r.SampleTime,
			r.BPM,
			r.HrvRMSSD,
			r.HrvSDNN,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
