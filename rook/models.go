package rook

type StepsResponse struct {
	Version    int    `json:"version"`
	ClientUUID string `json:"client_uuid"`
	UserID     string `json:"user_id"`

	PhysicalHealth struct {
		Summary struct {
			PhysicalSummary struct {
				Distance StepDistance    `json:"distance"`
				Metadata SummaryMetadata `json:"metadata"`
			} `json:"physical_summary"`
		} `json:"summary"`
	} `json:"physical_health"`
}

type StepDistance struct {
	Steps int `json:"steps_int"`
}

type SummaryMetadata struct {
	DateTime           string   `json:"datetime_string"`
	SourcesOfDataArray []string `json:"sources_of_data_array"`
}

// GlucoseResponse captures the body health blood glucose events payload.
type GlucoseResponse struct {
	Version       int    `json:"version"`
	DataStructure string `json:"data_structure"`
	ClientUUID    string `json:"client_uuid"`
	UserID        string `json:"user_id"`

	BodyHealth struct {
		Events struct {
			BloodGlucoseEvent []GlucoseEvent `json:"blood_glucose_event"`
		} `json:"events"`
	} `json:"body_health"`
}

type GlucoseEvent struct {
	Metadata     SummaryMetadata `json:"metadata"`
	BloodGlucose struct {
		AverageMgPerDL *int            `json:"blood_glucose_avg_mg_per_dL_int"`
		Samples        []GlucoseSample `json:"blood_glucose_granular_data_array"`
	} `json:"blood_glucose"`
}

type GlucoseSample struct {
	ValueMgPerDL int    `json:"blood_glucose_mg_per_dL_int"`
	DateTime     string `json:"datetime_string"`
}

type HeartRateResponse struct {
	Version    int    `json:"version"`
	ClientUUID string `json:"client_uuid"`
	UserID     string `json:"user_id"`

	PhysicalHealth struct {
		Events struct {
			HeartRateEvent []HeartRateEvent `json:"heart_rate_event"`
		} `json:"events"`
	} `json:"physical_health"`
}

type HeartRateEvent struct {
	Metadata  SummaryMetadata `json:"metadata"`
	HeartRate HeartRateData   `json:"heart_rate"`
}

type HeartRateData struct {
	AverageBPM      *int              `json:"hr_avg_bpm_int"`
	MinBPM          *int              `json:"hr_minimum_bpm_int"`
	MaxBPM          *int              `json:"hr_maximum_bpm_int"`
	RestingBPM      *int              `json:"hr_resting_bpm_int"`
	BPMSamples      []HeartRateSample `json:"hr_granular_data_array"`
	HrvRMSSDSamples []HRVRmssdSample  `json:"hrv_rmssd_granular_data_array"`
	HrvSDNNSamples  []HRVSdnnSample   `json:"hrv_sdnn_granular_data_array"`
	HrvAvgRMSSD     *float64          `json:"hrv_avg_rmssd_float"`
	HrvAvgSDNN      *float64          `json:"hrv_avg_sdnn_float"`
}

type HeartRateSample struct {
	BPM      int     `json:"hr_bpm_int"`
	DateTime string  `json:"datetime_string"`
	Interval float64 `json:"interval_duration_seconds_float"`
}

type HRVRmssdSample struct {
	Value    float64 `json:"hrv_rmssd_float"`
	DateTime string  `json:"datetime_string"`
	Interval float64 `json:"interval_duration_seconds_float"`
}

type HRVSdnnSample struct {
	Value    float64 `json:"hrv_sdnn_float"`
	DateTime string  `json:"datetime_string"`
	Interval float64 `json:"interval_duration_seconds_float"`
}
