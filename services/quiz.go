package services

import (
	"ego/models"
	"ego/repositories"
	"math"
)

// ──────────────────────────────────────────────────────────────
// Question Definition — metadata setiap soal IQ Test
// ──────────────────────────────────────────────────────────────

type questionDef struct {
	ID            string  // e.g., "Q_LR_001"
	Dikotomi      string  // "LR" | "NA" | "SA" | "LV"
	PolePrimary   string  // "L"|"R"|"N"|"A"|"S"|"A"|"L"|"V"
	Weight        float64 // bobot soal
	ReverseScored bool    // apakah reverse scored
}

// questionBank adalah bank soal IQ Test (20 soal)
var questionBank = []questionDef{
	// L/R (5 soal)
	{ID: "Q_LR_001", Dikotomi: "LR", PolePrimary: "L", Weight: 2.0, ReverseScored: false},
	{ID: "Q_LR_002", Dikotomi: "LR", PolePrimary: "R", Weight: 2.0, ReverseScored: false},
	{ID: "Q_LR_003", Dikotomi: "LR", PolePrimary: "L", Weight: 1.5, ReverseScored: false},
	{ID: "Q_LR_004", Dikotomi: "LR", PolePrimary: "R", Weight: 1.5, ReverseScored: false},
	{ID: "Q_LR_005", Dikotomi: "LR", PolePrimary: "L", Weight: 1.5, ReverseScored: true},

	// N/A (6 soal)
	{ID: "Q_NA_001", Dikotomi: "NA", PolePrimary: "N", Weight: 2.0, ReverseScored: false},
	{ID: "Q_NA_002", Dikotomi: "NA", PolePrimary: "A", Weight: 2.0, ReverseScored: false},
	{ID: "Q_NA_003", Dikotomi: "NA", PolePrimary: "N", Weight: 1.5, ReverseScored: false},
	{ID: "Q_NA_004", Dikotomi: "NA", PolePrimary: "A", Weight: 1.5, ReverseScored: false},
	{ID: "Q_NA_005", Dikotomi: "NA", PolePrimary: "N", Weight: 1.5, ReverseScored: true},
	{ID: "Q_NA_006", Dikotomi: "NA", PolePrimary: "A", Weight: 1.5, ReverseScored: true},

	// S/A (5 soal)
	{ID: "Q_SA_001", Dikotomi: "SA", PolePrimary: "S", Weight: 2.0, ReverseScored: false},
	{ID: "Q_SA_002", Dikotomi: "SA", PolePrimary: "A", Weight: 2.0, ReverseScored: false},
	{ID: "Q_SA_003", Dikotomi: "SA", PolePrimary: "S", Weight: 1.5, ReverseScored: false},
	{ID: "Q_SA_004", Dikotomi: "SA", PolePrimary: "A", Weight: 1.5, ReverseScored: false},
	{ID: "Q_SA_005", Dikotomi: "SA", PolePrimary: "S", Weight: 1.5, ReverseScored: true},

	// L/V (4 soal)
	{ID: "Q_LV_001", Dikotomi: "LV", PolePrimary: "L", Weight: 2.0, ReverseScored: false},
	{ID: "Q_LV_002", Dikotomi: "LV", PolePrimary: "V", Weight: 2.0, ReverseScored: false},
	{ID: "Q_LV_003", Dikotomi: "LV", PolePrimary: "L", Weight: 1.5, ReverseScored: false},
	{ID: "Q_LV_004", Dikotomi: "LV", PolePrimary: "V", Weight: 1.5, ReverseScored: true},
}

// ──────────────────────────────────────────────────────────────
// Likert contribution mapping (1–6)
// ──────────────────────────────────────────────────────────────

var likertContribution = map[int]float64{
	1: 1.00, // Sangat kuat ke pole_primary
	2: 0.67, // Kuat ke pole_primary
	3: 0.33, // Lemah ke pole_primary
	4: 0.33, // Lemah ke pole_opposite
	5: 0.67, // Kuat ke pole_opposite
	6: 1.00, // Sangat kuat ke pole_opposite
}

// ──────────────────────────────────────────────────────────────
// DeriveCognitiveProfile — menurunkan profil kognitif dari 4 huruf IQ Test type
// ──────────────────────────────────────────────────────────────

func DeriveCognitiveProfile(iqType string) models.CognitiveProfile {
	if len(iqType) < 4 {
		return models.CognitiveProfile{}
	}

	lR := string(iqType[0]) // "L" atau "R"
	nA := string(iqType[1]) // "N" atau "A"
	sA := string(iqType[2]) // "S" atau "A"
	lV := string(iqType[3]) // "L" atau "V"

	var dominant, auxiliary, complementary, developing string

	// Derivation algorithm per IQTEST.md §3.4
	if lR == "L" {
		if nA == "N" {
			dominant = "Logical"
			auxiliary = "Numerical"
		} else {
			dominant = "Logical"
			auxiliary = "Spatial"
		}
	} else {
		if sA == "S" {
			dominant = "Spatial"
			auxiliary = "Reasoning"
		} else {
			dominant = "Numerical"
			auxiliary = "Verbal"
		}
	}

	// Complementary and Developing are derived from the remaining dimensions
	// Simplified production rules based on IQTEST.md §3.4
	if lR == "L" && nA == "N" && sA == "S" && lV == "V" {
		complementary = "Spatial"
		developing = "Verbal"
	} else if lR == "L" && nA == "A" && sA == "A" && lV == "L" {
		complementary = "Analytical"
		developing = "Verbal"
	} else if lR == "R" && nA == "N" && sA == "S" && lV == "V" {
		complementary = "Numerical"
		developing = "Verbal"
	} else if lR == "R" && nA == "A" && sA == "A" && lV == "V" {
		complementary = "Reasoning"
		developing = "Linguistic"
	} else {
		// Fallback: assign remaining poles
		if lV == "L" {
			complementary = "Linguistic"
			developing = "Verbal"
		} else {
			complementary = "Verbal"
			developing = "Linguistic"
		}
	}

	return models.CognitiveProfile{
		Dominant:      dominant,
		Auxiliary:     auxiliary,
		Complementary: complementary,
		Developing:    developing,
	}
}

// ──────────────────────────────────────────────────────────────
// buildDimensionScore — menghitung DimensionScore dari akumulator
// ──────────────────────────────────────────────────────────────

func buildDimensionScore(poleA, poleB, max float64, poleALetter, poleBLetter string) models.DimensionScore {
	rawScore := poleA - poleB
	preference := poleALetter
	if rawScore < 0 {
		preference = poleBLetter
	}

	sci := 0.0
	if max > 0 {
		sci = math.Abs(rawScore) / max * 100
	}
	sci = math.Round(sci*10) / 10

	strength := "very_clear"
	switch {
	case sci <= 25:
		strength = "slight"
	case sci <= 50:
		strength = "moderate"
	case sci <= 75:
		strength = "clear"
	}

	return models.DimensionScore{
		RawScore:    rawScore,
		PoleAScore:  poleA,
		PoleBScore:  poleB,
		MaxPossible: max,
		Preference:  preference,
		SCI:         sci,
		Strength:    strength,
	}
}

// ──────────────────────────────────────────────────────────────
// CalculateIQResult — menghitung hasil IQ Test dari jawaban
// ──────────────────────────────────────────────────────────────

func CalculateIQResult(answers map[string]float64) models.IQTestResult {
	// Inisialisasi akumulator per dimensi
	type acc struct {
		poleA float64
		poleB float64
		max   float64
	}

	accumulators := map[string]*acc{
		"LR": {},
		"NA": {},
		"SA": {},
		"LV": {},
	}

	// Proses setiap jawaban
	for _, q := range questionBank {
		answerValue, ok := answers[q.ID]
		if !ok {
			continue
		}

		acc := accumulators[q.Dikotomi]

		// Skala Likert 1–6
		raw := int(answerValue)
		adjusted := raw
		if q.ReverseScored {
			adjusted = 7 - raw
		}

		contribution := likertContribution[adjusted]
		weighted := contribution * q.Weight

		// Tentukan pole: pole_primary = A, pole_opposite = B
		isPoleA := false
		switch q.Dikotomi {
		case "LR":
			isPoleA = q.PolePrimary == "L"
		case "NA":
			isPoleA = q.PolePrimary == "N"
		case "SA":
			isPoleA = q.PolePrimary == "S"
		case "LV":
			isPoleA = q.PolePrimary == "L"
		}

		if adjusted <= 3 {
			// Condong ke pole_primary
			if isPoleA {
				acc.poleA += weighted
			} else {
				acc.poleB += weighted
			}
		} else {
			// Condong ke pole_opposite
			if isPoleA {
				acc.poleB += weighted
			} else {
				acc.poleA += weighted
			}
		}

		acc.max += q.Weight
	}

	// Hitung DimensionScore untuk setiap dimensi
	scores := map[string]models.DimensionScore{
		"LR": buildDimensionScore(accumulators["LR"].poleA, accumulators["LR"].poleB, accumulators["LR"].max, "L", "R"),
		"NA": buildDimensionScore(accumulators["NA"].poleA, accumulators["NA"].poleB, accumulators["NA"].max, "N", "A"),
		"SA": buildDimensionScore(accumulators["SA"].poleA, accumulators["SA"].poleB, accumulators["SA"].max, "S", "A"),
		"LV": buildDimensionScore(accumulators["LV"].poleA, accumulators["LV"].poleB, accumulators["LV"].max, "L", "V"),
	}

	// Derive IQ Test type
	iqType := scores["LR"].Preference +
		scores["NA"].Preference +
		scores["SA"].Preference +
		scores["LV"].Preference

	// Derive cognitive profile
	cognitiveProfile := DeriveCognitiveProfile(iqType)

	return models.IQTestResult{
		Type:             iqType,
		Scores:           scores,
		CognitiveProfile: cognitiveProfile,
	}
}

// ──────────────────────────────────────────────────────────────
// ProcessQuizAnswers — memproses 20 jawaban kuesioner IQ Test
// ──────────────────────────────────────────────────────────────

func ProcessQuizAnswers(email, nama string, rawAnswers map[string]float64) (string, error) {
	// Hitung IQ Test
	result := CalculateIQResult(rawAnswers)

	// Ambil raw scores (integer) untuk disimpan di database
	skorLR := int(result.Scores["LR"].RawScore)
	skorNA := int(result.Scores["NA"].RawScore)
	skorSA := int(result.Scores["SA"].RawScore)
	skorLV := int(result.Scores["LV"].RawScore)

	userID, err := repositories.InsertUser(email, nama, skorLR, skorNA, skorSA, skorLV, result.Type)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// ──────────────────────────────────────────────────────────────
// GetPaywallData — mengambil data untuk halaman paywall
// ──────────────────────────────────────────────────────────────

func GetPaywallData(id string) (*models.PaywallData, error) {
	nama, err := repositories.GetUserName(id)
	if err != nil {
		return nil, err
	}
	return &models.PaywallData{ID: id, Nama: nama}, nil
}

// ──────────────────────────────────────────────────────────────
// GetQuizResult — mengambil data hasil kuis (dengan proteksi paywall)
// ──────────────────────────────────────────────────────────────

func mapIQToDarkTriad(skorLR, skorNA, skorSA, _ int) (narsisme, machiavellian, psikopati int) {
	// Map IQ Test raw scores to Dark Triad percentile-like values (0-100)
	// Per IQTEST.md §8.3: L/R → Narcissism, N/A → Machiavellianism, S/A → Psychopathy
	// Use absolute values capped at a reasonable scale
	narsisme = absInt(skorLR) * 5 // L/R dimension → Narsisme
	if narsisme > 100 {
		narsisme = 100
	}
	machiavellian = absInt(skorNA) * 5 // N/A dimension → Machiavellian
	if machiavellian > 100 {
		machiavellian = 100
	}
	psikopati = absInt(skorSA) * 5 // S/A dimension → Psikopati
	if psikopati > 100 {
		psikopati = 100
	}
	return
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GetQuizResult(id string) (*models.QuizResult, error) {
	user, err := repositories.GetUserResult(id)
	if err != nil {
		return nil, err
	}

	// Proteksi: hanya tampilkan hasil jika sudah PAID
	if user.StatusPembayaran != "PAID" {
		return nil, nil // nil menandakan belum bayar
	}

	// Build scores map for template
	scores := map[string]models.DimensionScore{
		"LR": buildDimensionScore(float64(user.SkorLR), 0, float64(user.SkorLR), "L", "R"),
		"NA": buildDimensionScore(float64(user.SkorNA), 0, float64(user.SkorNA), "N", "A"),
		"SA": buildDimensionScore(float64(user.SkorSA), 0, float64(user.SkorSA), "S", "A"),
		"LV": buildDimensionScore(float64(user.SkorLV), 0, float64(user.SkorLV), "L", "V"),
	}

	// Map IQ Test raw scores to Dark Triad dimensions for narrative generation
	narsisme, machiavellian, psikopati := mapIQToDarkTriad(user.SkorLR, user.SkorNA, user.SkorSA, user.SkorLV)

	// Generate all narratives using the Dark Triad scoring system
	execSummary, relProfile, kekuatan, areaPerhatian, relInsight, compatNotes, refQuestions :=
		GenerateAllNarratives(user.Nama, narsisme, machiavellian, psikopati)

	return &models.QuizResult{
		Nama:                user.Nama,
		IQTipe:              user.IQTipe,
		SkorLR:              user.SkorLR,
		SkorNA:              user.SkorNA,
		SkorSA:              user.SkorSA,
		SkorLV:              user.SkorLV,
		Scores:              scores,
		CognitiveProfile:    DeriveCognitiveProfile(user.IQTipe),
		ExecutiveSummary:    execSummary,
		RelationshipProfile: relProfile,
		Kekuatan:            kekuatan,
		AreaPerhatian:       areaPerhatian,
		RelationshipInsight: relInsight,
		CompatibilityNotes:  compatNotes,
		ReflectionQuestions: refQuestions,
	}, nil
}

// ──────────────────────────────────────────────────────────────
// ConfirmPayment — mengonfirmasi pembayaran user
// ──────────────────────────────────────────────────────────────

func ConfirmPayment(id string) error {
	return repositories.UpdatePaymentStatus(id)
}

// ──────────────────────────────────────────────────────────────
// GetAllUsers — mengambil semua user untuk admin
// ──────────────────────────────────────────────────────────────

func GetAllUsers() ([]models.User, error) {
	return repositories.GetAllUsers()
}

// ──────────────────────────────────────────────────────────────
// GetUserByID — mengambil user by ID untuk admin
// ──────────────────────────────────────────────────────────────

func GetUserByID(id string) (*models.User, error) {
	return repositories.GetUserByID(id)
}
