package models

import "time"

// User — data pengguna (per IQTEST.md §10.2 users table)
// Score fields dihapus — skor disimpan di iq_results
type User struct {
	ID        string
	Email     string
	Nama      string
	Phone     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// QuizResult — data yang dikirim ke template hasil (per IQTEST.md §8.2)
type QuizResult struct {
	ID   string
	Nama string

	// Skor
	RawScore    float64
	MaxPossible float64
	Percentile  *float64
	EstimatedIQ *float64

	// Skor per domain
	DomainScores map[string]DomainScore

	// Waktu
	AvgResponseMs int

	// Reliabilitas
	IsReliable       bool
	ReliabilityFlags []string

	// Narrative fields
	ExecutiveSummary string
	Kekuatan         []string
	AreaPerhatian    []string
}

// PaywallData — data untuk rendering paywall page
type PaywallData struct {
	ID   string
	Nama string
}

// Payment — data pembayaran (per IQTEST.md §10.2 payments table)
type Payment struct {
	ID            string
	UserID        string
	SessionID     string
	Amount        float64
	Currency      string
	Status        string
	PaymentMethod *string
	PaidAt        *interface{}
	CreatedAt     interface{}
}
