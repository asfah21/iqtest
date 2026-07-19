package models

import "time"

// SessionResponse — jawaban user untuk satu soal
type SessionResponse struct {
	QuestionID     string
	QuestionCode   string
	SelectedOption *string // "A"/"B"/"C"/"D" atau nil jika timeout
	IsCorrect      bool
	TimeTakenMs    int
	TimedOut       bool
}

// Session — data sesi tes (per IQTEST.md §10.2 test_sessions table)
type Session struct {
	ID           string
	UserID       string
	SessionToken string
	StartedAt    time.Time
	CompletedAt  *time.Time
	DeviceType   string
	IPAddress    string
	IsCompleted  bool
	Metadata     interface{} // JSONB
}
