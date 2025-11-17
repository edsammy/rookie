package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"rookie/db"
	"rookie/rook"
	"time"

	"github.com/joho/godotenv"
)

const userID = "david1"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	db.Init()
	defer db.DB.Close()

	ctx := context.Background()

	end := time.Now()
	start := end.AddDate(0, 0, -7)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		resp, err := rook.FetchBloodGlucose(ctx, userID, d)
		if err != nil {
			fmt.Printf("glucose fetch %s: %v\n", d.Format("2006-01-02"), err)
			continue
		}
		if resp == nil {
			fmt.Printf("glucose fetch %s: no data\n", d.Format("2006-01-02"))
			continue
		}

		var readings []db.BloodGlucoseReading
		for _, event := range resp.BodyHealth.Events.BloodGlucoseEvent {
			for _, sample := range event.BloodGlucose.Samples {
				sampleTime, err := time.Parse(time.RFC3339Nano, sample.DateTime)
				if err != nil {
					fmt.Printf("glucose parse %s: %v\n", sample.DateTime, err)
					continue
				}

				readings = append(readings, db.BloodGlucoseReading{
					UserID:     userID,
					SampleTime: sampleTime,
					MgDL:       sql.NullInt64{Int64: int64(sample.ValueMgPerDL), Valid: true},
				})
			}
		}

		if err := db.UpsertBloodGlucoseBatch(ctx, readings); err != nil {
			fmt.Printf("glucose insert %s: %v\n", d.Format("2006-01-02"), err)
		}
	}

	end = time.Now()
	start = end.AddDate(0, 0, -29)

	var stepTotals []db.StepTotal
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		resp, err := rook.FetchSteps(ctx, userID, d)
		if err != nil {
			fmt.Printf("steps fetch %s: %v\n", d.Format("2006-01-02"), err)
			continue
		}
		if resp == nil {
			fmt.Printf("steps fetch %s: no data\n", d.Format("2006-01-02"))
			continue
		}

		steps := resp.PhysicalHealth.Summary.PhysicalSummary.Distance.Steps
		stepTotals = append(stepTotals, db.StepTotal{
			UserID: userID,
			Date:   d,
			Total:  steps,
		})
	}

	if err := db.UpsertStepsBatch(ctx, stepTotals); err != nil {
		fmt.Printf("steps insert: %v\n", err)
	}

	end = time.Now()
	start = end.AddDate(0, 0, -7)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		fmt.Printf("getting hr data")
		resp, err := rook.FetchHeartRate(ctx, userID, d)
		if err != nil {
			fmt.Printf("hr fetch %s: %v\n", d.Format("2006-01-02"), err)
			continue
		}
		if resp == nil {
			fmt.Printf("hr fetch %s: no data\n", d.Format("2006-01-02"))
			continue
		}

		raw, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Printf("marshal heart rate %s: %v\n", d.Format("2006-01-02"), err)
		} else {
			fmt.Println(string(raw))
		}

		readingsMap := make(map[time.Time]*db.HeartRateReading)

		record := func(ts time.Time) *db.HeartRateReading {
			if r, ok := readingsMap[ts]; ok {
				return r
			}
			r := &db.HeartRateReading{
				UserID:     userID,
				SampleTime: ts,
			}
			readingsMap[ts] = r
			return r
		}

		for _, event := range resp.PhysicalHealth.Events.HeartRateEvent {
			for _, sample := range event.HeartRate.BPMSamples {
				ts, err := time.Parse(time.RFC3339Nano, sample.DateTime)
				if err != nil {
					fmt.Printf("hr bpm parse %s: %v\n", sample.DateTime, err)
					continue
				}
				r := record(ts)
				r.BPM = sql.NullInt64{Int64: int64(sample.BPM), Valid: true}
			}

			for _, sample := range event.HeartRate.HrvRMSSDSamples {
				ts, err := time.Parse(time.RFC3339Nano, sample.DateTime)
				if err != nil {
					fmt.Printf("hr rmssd parse %s: %v\n", sample.DateTime, err)
					continue
				}
				r := record(ts)
				r.HrvRMSSD = sql.NullFloat64{Float64: sample.Value, Valid: true}
			}

			for _, sample := range event.HeartRate.HrvSDNNSamples {
				ts, err := time.Parse(time.RFC3339Nano, sample.DateTime)
				if err != nil {
					fmt.Printf("hr sdnn parse %s: %v\n", sample.DateTime, err)
					continue
				}
				r := record(ts)
				r.HrvSDNN = sql.NullFloat64{Float64: sample.Value, Valid: true}
			}
		}

		var readings []db.HeartRateReading
		for _, r := range readingsMap {
			readings = append(readings, *r)
		}

		if err := db.UpsertHeartRateBatch(ctx, readings); err != nil {
			fmt.Printf("hr insert %s: %v\n", d.Format("2006-01-02"), err)
		}
	}
}
