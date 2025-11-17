package db

import (
	"context"
	"time"
)

type StepTotal struct {
	UserID string
	Date   time.Time
	Total  int
}

const upsertSteps = `
INSERT INTO steps (user_id, activity_date, total_steps)
VALUES (?, ?, ?)
ON CONFLICT(user_id, activity_date) DO UPDATE SET
	total_steps = excluded.total_steps;
`

func UpsertSteps(ctx context.Context, total StepTotal) error {
	_, err := DB.ExecContext(ctx, upsertSteps, total.UserID, total.Date.Format("2006-01-02"), total.Total)
	return err
}

func UpsertStepsBatch(ctx context.Context, totals []StepTotal) error {
	if len(totals) == 0 {
		return nil
	}

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, upsertSteps)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, t := range totals {
		if _, err := stmt.ExecContext(ctx, t.UserID, t.Date.Format("2006-01-02"), t.Total); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
