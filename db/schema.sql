CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	name TEXT,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS steps (
	user_id TEXT NOT NULL,
	activity_date DATE NOT NULL,
	total_steps INTEGER NOT NULL DEFAULT 0,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, activity_date)
);

CREATE TABLE IF NOT EXISTS blood_glucose (
	user_id TEXT NOT NULL,
	sample_time DATETIME NOT NULL,
	mg_dl INTEGER,
	PRIMARY KEY (user_id, sample_time)
);

CREATE TABLE IF NOT EXISTS heart_rate (
	user_id TEXT NOT NULL,
	sample_time DATETIME NOT NULL,
	bpm INTEGER,
	hrv_rmssd REAL,
	hrv_sdnn REAL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, sample_time)
);
