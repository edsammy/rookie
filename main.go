package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rookie/db"
	"time"
)

func main() {
	db.Init()
	defer db.DB.Close()

	mux := http.NewServeMux()
	mux.Handle("/",
		http.StripPrefix("/",
			http.FileServer(http.Dir("web"))))
	mux.HandleFunc("/api/glucose", glucoseAPIHandler)
	mux.HandleFunc("/api/steps", stepsAPIHandler)
	mux.HandleFunc("/api/heart_rate", heartRateAPIHandler)

	addr := ":8080"
	log.Printf("serving data on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

type glucosePoint struct {
	SampleTime time.Time `json:"sample_time"`
	MgDL       *int      `json:"mg_dl"`
}

func glucoseAPIHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	start, end, err := resolveRange(r, 7)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const query = `
		SELECT sample_time, mg_dl
		FROM blood_glucose
		WHERE user_id = ?
		  AND sample_time BETWEEN ? AND ?
		ORDER BY sample_time
	`

	rows, err := db.DB.QueryContext(r.Context(), query, userID, start, end)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var points []glucosePoint
	for rows.Next() {
		var (
			pt  glucosePoint
			mg  sql.NullInt64
			err = rows.Scan(&pt.SampleTime, &mg)
		)
		if err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		if mg.Valid {
			val := int(mg.Int64)
			pt.MgDL = &val
		}
		points = append(points, pt)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "rows error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(points); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

type heartRatePoint struct {
	SampleTime time.Time `json:"sample_time"`
	BPM        *int      `json:"bpm,omitempty"`
	HrvRMSSD   *float64  `json:"hrv_rmssd,omitempty"`
	HrvSDNN    *float64  `json:"hrv_sdnn,omitempty"`
}

type stepsPoint struct {
	ActivityDate time.Time `json:"activity_date"`
	TotalSteps   int       `json:"total_steps"`
}

func heartRateAPIHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	start, end, err := resolveRange(r, 7)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const query = `
		SELECT sample_time, bpm, hrv_rmssd, hrv_sdnn
		FROM heart_rate
		WHERE user_id = ?
		  AND sample_time BETWEEN ? AND ?
		ORDER BY sample_time
	`

	rows, err := db.DB.QueryContext(r.Context(), query, userID, start, end)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var points []heartRatePoint
	for rows.Next() {
		var (
			pt    heartRatePoint
			bpm   sql.NullInt64
			rmssd sql.NullFloat64
			sdnn  sql.NullFloat64
		)
		if err := rows.Scan(&pt.SampleTime, &bpm, &rmssd, &sdnn); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		if bpm.Valid {
			val := int(bpm.Int64)
			pt.BPM = &val
		}
		if rmssd.Valid {
			val := rmssd.Float64
			pt.HrvRMSSD = &val
		}
		if sdnn.Valid {
			val := sdnn.Float64
			pt.HrvSDNN = &val
		}
		points = append(points, pt)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "rows error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(points); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

func stepsAPIHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	start, end, err := resolveRange(r, 30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	const query = `
		SELECT activity_date, total_steps
		FROM steps
		WHERE user_id = ?
		  AND activity_date BETWEEN ? AND ?
		ORDER BY activity_date
	`

	rows, err := db.DB.QueryContext(r.Context(), query, userID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var points []stepsPoint
	for rows.Next() {
		var (
			pt  stepsPoint
			err = rows.Scan(&pt.ActivityDate, &pt.TotalSteps)
		)
		if err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		points = append(points, pt)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "rows error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(points); err != nil {
		http.Error(w, "encode error", http.StatusInternalServerError)
		return
	}
}

func parseDate(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("bad date")
}

func resolveRange(r *http.Request, defaultDays int) (time.Time, time.Time, error) {
	end := time.Now()
	if endStr := r.URL.Query().Get("end"); endStr != "" {
		t, err := parseDate(endStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end")
		}
		end = t
	}

	start := end.AddDate(0, 0, -defaultDays)
	if startStr := r.URL.Query().Get("start"); startStr != "" {
		t, err := parseDate(startStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start")
		}
		start = t
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("start must be before end")
	}

	return start, end, nil
}
