package models

// DomainScore — skor untuk satu domain kognitif
type DomainScore struct {
	Domain      string  // "MTX" | "SEQ" | "SPA" | "ANL"
	RawScore    float64 // skor tertimbang
	MaxPossible float64 // skor maksimum domain ini
	Percentage  float64 // (raw/max) * 100
}

// IQTestResult — output akhir kalkulasi (per IQTEST.md §6)
type IQTestResult struct {
	RawScore         float64
	MaxPossible      float64
	DomainScores     map[string]DomainScore
	Percentile       *float64 // NULL sampai data normatif tersedia
	EstimatedIQ      *float64 // NULL sampai data normatif tersedia
	AvgResponseMs    int
	IsReliable       bool
	ReliabilityFlags []string
}

// ReliabilityFlag — indikator keandalan hasil tes (per IQTEST.md §9.3)
type ReliabilityFlag struct {
	IsReliable     bool
	Reasons        []string // "speed_guessing", "tab_switch_excessive", dll
	Recommendation string   // "hasil_valid" | "disarankan_mengulang"
}
