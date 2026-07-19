package models

// User merepresentasikan data pengguna yang mengisi kuesioner MBTI
type User struct {
	ID               string `json:"id"`
	Nama             string `json:"nama"`
	Email            string `json:"email"`
	SkorEI           int    `json:"skor_ei"`   // Raw score E/I (positif = E, negatif = I)
	SkorSN           int    `json:"skor_sn"`   // Raw score S/N (positif = S, negatif = N)
	SkorTF           int    `json:"skor_tf"`   // Raw score T/F (positif = T, negatif = F)
	SkorJP           int    `json:"skor_jp"`   // Raw score J/P (positif = J, negatif = P)
	MBTITipe         string `json:"mbti_tipe"` // e.g., "INTJ"
	StatusPembayaran string `json:"status_pembayaran"`
}

// DikotomiScore menyimpan hasil skoring untuk satu sumbu dikotomi
type DikotomiScore struct {
	RawScore    float64 `json:"raw_score"`
	PoleAScore  float64 `json:"pole_a_score"`
	PoleBScore  float64 `json:"pole_b_score"`
	MaxPossible float64 `json:"max_possible"`
	Preference  string  `json:"preference"`
	PCI         float64 `json:"pci"`
	Strength    string  `json:"strength"`
}

// CognitiveStack menyimpan urutan 4 fungsi kognitif hasil derivasi
type CognitiveStack struct {
	Dominant  string `json:"dominant"`
	Auxiliary string `json:"auxiliary"`
	Tertiary  string `json:"tertiary"`
	Inferior  string `json:"inferior"`
}

// MBTIResult adalah output akhir kalkulasi satu sesi tes
type MBTIResult struct {
	Type           string                   `json:"type"`
	Scores         map[string]DikotomiScore `json:"scores"`
	CognitiveStack CognitiveStack           `json:"cognitive_stack"`
}

// QuizResult adalah data yang dikirim ke template hasil
type QuizResult struct {
	Nama string
	MBTI string

	// Raw scores
	SkorEI int
	SkorSN int
	SkorTF int
	SkorJP int

	// Dikotomi scores
	Scores map[string]DikotomiScore

	// Cognitive stack
	CognitiveStack CognitiveStack

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
