package rook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	physicalBase = "https://api.rook-connect.review/v2/processed_data/physical_health"
	bodyBase     = "https://api.rook-connect.review/v2/processed_data/body_health"
)

func FetchSteps(ctx context.Context, userID string, date time.Time) (*StepsResponse, error) {
	resp, err := doRequest(ctx, physicalBase, "/summary", map[string]string{
		"user_id": userID,
		"date":    date.Format("2006-01-02"),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rook steps: %s", resp.Status)
	}

	var payload StepsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func FetchBloodGlucose(ctx context.Context, userID string, date time.Time) (*GlucoseResponse, error) {
	resp, err := doRequest(ctx, bodyBase, "/events/blood_glucose", map[string]string{
		"user_id": userID,
		"date":    date.Format("2006-01-02"),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rook glucose: %s", resp.Status)
	}

	var payload GlucoseResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func FetchHeartRate(ctx context.Context, userID string, date time.Time) (*HeartRateResponse, error) {
	resp, err := doRequest(ctx, physicalBase, "/events/heart_rate", map[string]string{
		"user_id": userID,
		"date":    date.Format("2006-01-02"),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rook heart rate: %s", resp.Status)
	}

	var payload HeartRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func doRequest(ctx context.Context, base, path string, params map[string]string) (*http.Response, error) {
	fullURL, err := url.Parse(base + path)
	if err != nil {
		return nil, err
	}

	q := fullURL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	fullURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("ROOK_CLIENT_UUID"), os.Getenv("ROOK_SECRET_KEY"))
	return http.DefaultClient.Do(req)
}
