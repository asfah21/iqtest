package models

// User merepresentasikan data pengguna yang mengisi kuesioner IQ Test
type User struct {
	ID               string `json:"id"`
	Nama             string `json:"nama"`
	Email            string `json:"email"`
	SkorLR           int    `json:"skor_lr"` // Raw score L/R (positif = L, negatif = R)
	SkorNA           int    `json:"skor_na"` // Raw score N/A (positif = N, negatif = A)
	SkorSA           int    `json:"skor_sa"` // Raw score S/A (positif = S, negatif = A)
	SkorLV           int    `json:"skor_lv"` // Raw score L/V (positif = L, negatif = V)
	IQTipe           string `json:"iq_tipe"` // e.g., "LNSL"
	StatusPembayaran string `json:"status_pembayaran"`
}

// DimensionScore menyimpan hasil skoring untuk satu dimensi kognitif
type DimensionScore struct {
	RawScore    float64 `json:"raw_score"`
	PoleAScore  float64 `json:"pole_a_score"`
	PoleBScore  float64 `json:"pole_b_score"`
	MaxPossible float64 `json:"max_possible"`
	Preference  string  `json:"preference"`
	SCI         float64 `json:"sci"`
	Strength    string  `json:"strength"`
}

// CognitiveProfile menyimpan urutan 4 kemampuan kognitif hasil derivasi
type CognitiveProfile struct {
	Dominant      string `json:"dominant"`
	Auxiliary     string `json:"auxiliary"`
	Complementary string `json:"complementary"`
	Developing    string `json:"developing"`
}

// IQTestResult adalah output akhir kalkulasi satu sesi tes
type IQTestResult struct {
	Type             string                    `json:"type"`
	Scores           map[string]DimensionScore `json:"scores"`
	CognitiveProfile CognitiveProfile          `json:"cognitive_profile"`
}

// QuizResult adalah data yang dikirim ke template hasil
type QuizResult struct {
	Nama   string
	IQTipe string

	// Raw scores
	SkorLR int
	SkorNA int
	SkorSA int
	SkorLV int

	// Dimension scores
	Scores map[string]DimensionScore

	// Cognitive profile
	CognitiveProfile CognitiveProfile

	// Narrative fields
	ExecutiveSummary    string
	RelationshipProfile string
	Kekuatan            []string
	AreaPerhatian       []string
	RelationshipInsight string
	CompatibilityNotes  string
	ReflectionQuestions []string
}

// PaywallData adalah data yang dikirim ke template paywall
type PaywallData struct {
	ID   string
	Nama string
}
