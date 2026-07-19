package services

import (
	"ego/models"
	"ego/repositories"
	"math"
)

// ──────────────────────────────────────────────────────────────
// Question Definition — metadata setiap soal MBTI
// ──────────────────────────────────────────────────────────────

type questionDef struct {
	ID            string  // e.g., "Q_EI_001"
	Dikotomi      string  // "EI" | "SN" | "TF" | "JP"
	PolePrimary   string  // "E"|"I"|"S"|"N"|"T"|"F"|"J"|"P"
	Weight        float64 // bobot soal
	ReverseScored bool    // apakah reverse scored
}

// questionBank adalah bank soal MBTI (20 soal)
var questionBank = []questionDef{
	// E/I (5 soal)
	{ID: "Q_EI_001", Dikotomi: "EI", PolePrimary: "E", Weight: 2.0, ReverseScored: false},
	{ID: "Q_EI_002", Dikotomi: "EI", PolePrimary: "I", Weight: 2.0, ReverseScored: false},
	{ID: "Q_EI_003", Dikotomi: "EI", PolePrimary: "E", Weight: 1.5, ReverseScored: false},
	{ID: "Q_EI_004", Dikotomi: "EI", PolePrimary: "I", Weight: 1.5, ReverseScored: false},
	{ID: "Q_EI_005", Dikotomi: "EI", PolePrimary: "E", Weight: 1.5, ReverseScored: true},

	// S/N (6 soal)
	{ID: "Q_SN_001", Dikotomi: "SN", PolePrimary: "S", Weight: 2.0, ReverseScored: false},
	{ID: "Q_SN_002", Dikotomi: "SN", PolePrimary: "N", Weight: 2.0, ReverseScored: false},
	{ID: "Q_SN_003", Dikotomi: "SN", PolePrimary: "S", Weight: 1.5, ReverseScored: false},
	{ID: "Q_SN_004", Dikotomi: "SN", PolePrimary: "N", Weight: 1.5, ReverseScored: false},
	{ID: "Q_SN_005", Dikotomi: "SN", PolePrimary: "S", Weight: 1.5, ReverseScored: true},
	{ID: "Q_SN_006", Dikotomi: "SN", PolePrimary: "N", Weight: 1.5, ReverseScored: true},

	// T/F (5 soal)
	{ID: "Q_TF_001", Dikotomi: "TF", PolePrimary: "T", Weight: 2.0, ReverseScored: false},
	{ID: "Q_TF_002", Dikotomi: "TF", PolePrimary: "F", Weight: 2.0, ReverseScored: false},
	{ID: "Q_TF_003", Dikotomi: "TF", PolePrimary: "T", Weight: 1.5, ReverseScored: false},
	{ID: "Q_TF_004", Dikotomi: "TF", PolePrimary: "F", Weight: 1.5, ReverseScored: false},
	{ID: "Q_TF_005", Dikotomi: "TF", PolePrimary: "T", Weight: 1.5, ReverseScored: true},

	// J/P (4 soal)
	{ID: "Q_JP_001", Dikotomi: "JP", PolePrimary: "J", Weight: 2.0, ReverseScored: false},
	{ID: "Q_JP_002", Dikotomi: "JP", PolePrimary: "P", Weight: 2.0, ReverseScored: false},
	{ID: "Q_JP_003", Dikotomi: "JP", PolePrimary: "J", Weight: 1.5, ReverseScored: false},
	{ID: "Q_JP_004", Dikotomi: "JP", PolePrimary: "P", Weight: 1.5, ReverseScored: true},
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
// Axis opposites untuk derivasi cognitive stack
// ──────────────────────────────────────────────────────────────

var axisOpposites = map[string]string{
	"Se": "Ni", "Ni": "Se",
	"Si": "Ne", "Ne": "Si",
	"Te": "Fi", "Fi": "Te",
	"Ti": "Fe", "Fe": "Ti",
}

func axisOpposite(fn string) string {
	return axisOpposites[fn]
}

// ──────────────────────────────────────────────────────────────
// DeriveCognitiveStack — menurunkan stack fungsi kognitif dari 4 huruf MBTI
// ──────────────────────────────────────────────────────────────

func DeriveCognitiveStack(mbti string) models.CognitiveStack {
	if len(mbti) < 4 {
		return models.CognitiveStack{}
	}

	eI := string(mbti[0]) // "E" atau "I"
	sN := string(mbti[1]) // "S" atau "N"
	tF := string(mbti[2]) // "T" atau "F"
	jP := string(mbti[3]) // "J" atau "P"

	// Tentukan fungsi persepsi
	perceiving := "Se"
	perceivingI := "Si"
	if sN == "N" {
		perceiving = "Ne"
		perceivingI = "Ni"
	}

	// Tentukan fungsi penilaian
	judging := "Te"
	judgingI := "Ti"
	if tF == "F" {
		judging = "Fe"
		judgingI = "Fi"
	}

	var dominant, auxiliary, tertiary, inferior string

	if eI == "E" {
		if jP == "J" {
			// Extravert-Judging: Dominant = fungsi judging ekstraverted
			dominant = judging
			auxiliary = perceivingI
			tertiary = axisOpposite(perceivingI)
			inferior = axisOpposite(dominant)
		} else {
			// Extravert-Perceiving: Dominant = fungsi perceiving ekstraverted
			dominant = perceiving
			auxiliary = judgingI
			tertiary = axisOpposite(judgingI)
			inferior = axisOpposite(dominant)
		}
	} else {
		if jP == "J" {
			// Introvert-Judging: Dominant = fungsi judging introverted
			dominant = judgingI
			auxiliary = perceiving
			tertiary = axisOpposite(perceiving)
			inferior = axisOpposite(dominant)
		} else {
			// Introvert-Perceiving: Dominant = fungsi perceiving introverted
			dominant = perceivingI
			auxiliary = judging
			tertiary = axisOpposite(judging)
			inferior = axisOpposite(dominant)
		}
	}

	return models.CognitiveStack{
		Dominant:  dominant,
		Auxiliary: auxiliary,
		Tertiary:  tertiary,
		Inferior:  inferior,
	}
}

// ──────────────────────────────────────────────────────────────
// buildDikotomiScore — menghitung DikotomiScore dari akumulator
// ──────────────────────────────────────────────────────────────

func buildDikotomiScore(poleA, poleB, max float64, poleALetter, poleBLetter string) models.DikotomiScore {
	rawScore := poleA - poleB
	preference := poleALetter
	if rawScore < 0 {
		preference = poleBLetter
	}

	pci := 0.0
	if max > 0 {
		pci = math.Abs(rawScore) / max * 100
	}
	pci = math.Round(pci*10) / 10

	strength := "very_clear"
	switch {
	case pci <= 25:
		strength = "slight"
	case pci <= 50:
		strength = "moderate"
	case pci <= 75:
		strength = "clear"
	}

	return models.DikotomiScore{
		RawScore:    rawScore,
		PoleAScore:  poleA,
		PoleBScore:  poleB,
		MaxPossible: max,
		Preference:  preference,
		PCI:         pci,
		Strength:    strength,
	}
}

// ──────────────────────────────────────────────────────────────
// CalculateMBTI — menghitung hasil MBTI dari jawaban
// ──────────────────────────────────────────────────────────────

func CalculateMBTI(answers map[string]float64) models.MBTIResult {
	// Inisialisasi akumulator per dikotomi
	type acc struct {
		poleA float64
		poleB float64
		max   float64
	}

	accumulators := map[string]*acc{
		"EI": {},
		"SN": {},
		"TF": {},
		"JP": {},
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
		case "EI":
			isPoleA = q.PolePrimary == "E"
		case "SN":
			isPoleA = q.PolePrimary == "S"
		case "TF":
			isPoleA = q.PolePrimary == "T"
		case "JP":
			isPoleA = q.PolePrimary == "J"
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

	// Hitung DikotomiScore untuk setiap dikotomi
	scores := map[string]models.DikotomiScore{
		"EI": buildDikotomiScore(accumulators["EI"].poleA, accumulators["EI"].poleB, accumulators["EI"].max, "E", "I"),
		"SN": buildDikotomiScore(accumulators["SN"].poleA, accumulators["SN"].poleB, accumulators["SN"].max, "S", "N"),
		"TF": buildDikotomiScore(accumulators["TF"].poleA, accumulators["TF"].poleB, accumulators["TF"].max, "T", "F"),
		"JP": buildDikotomiScore(accumulators["JP"].poleA, accumulators["JP"].poleB, accumulators["JP"].max, "J", "P"),
	}

	// Derive tipe MBTI
	mbtiType := scores["EI"].Preference +
		scores["SN"].Preference +
		scores["TF"].Preference +
		scores["JP"].Preference

	// Derive cognitive stack
	cognitiveStack := DeriveCognitiveStack(mbtiType)

	return models.MBTIResult{
		Type:           mbtiType,
		Scores:         scores,
		CognitiveStack: cognitiveStack,
	}
}

// ──────────────────────────────────────────────────────────────
// ProcessQuizAnswers — memproses 20 jawaban kuesioner MBTI
// ──────────────────────────────────────────────────────────────

func ProcessQuizAnswers(email, nama string, rawAnswers map[string]float64) (string, error) {
	// Hitung MBTI
	result := CalculateMBTI(rawAnswers)

	// Ambil raw scores (integer) untuk disimpan di database
	skorEI := int(result.Scores["EI"].RawScore)
	skorSN := int(result.Scores["SN"].RawScore)
	skorTF := int(result.Scores["TF"].RawScore)
	skorJP := int(result.Scores["JP"].RawScore)

	userID, err := repositories.InsertUser(email, nama, skorEI, skorSN, skorTF, skorJP, result.Type)
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

func mapMBTIToDarkTriad(skorEI, skorSN, skorTF, skorJP int) (narsisme, machiavellian, psikopati int) {
	// Map MBTI raw scores to Dark Triad percentile-like values (0-100)
	// Use absolute values capped at a reasonable scale
	narsisme = absInt(skorEI) * 5 // E/I dimension → Narsisme
	if narsisme > 100 {
		narsisme = 100
	}
	machiavellian = absInt(skorSN) * 5 // S/N dimension → Machiavellian
	if machiavellian > 100 {
		machiavellian = 100
	}
	psikopati = absInt(skorTF) * 5 // T/F dimension → Psikopati
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
	scores := map[string]models.DikotomiScore{
		"EI": buildDikotomiScore(float64(user.SkorEI), 0, float64(user.SkorEI), "E", "I"),
		"SN": buildDikotomiScore(float64(user.SkorSN), 0, float64(user.SkorSN), "S", "N"),
		"TF": buildDikotomiScore(float64(user.SkorTF), 0, float64(user.SkorTF), "T", "F"),
		"JP": buildDikotomiScore(float64(user.SkorJP), 0, float64(user.SkorJP), "J", "P"),
	}

	// Map MBTI raw scores to Dark Triad dimensions for narrative generation
	narsisme, machiavellian, psikopati := mapMBTIToDarkTriad(user.SkorEI, user.SkorSN, user.SkorTF, user.SkorJP)

	// Generate all narratives using the Dark Triad scoring system
	execSummary, relProfile, kekuatan, areaPerhatian, relInsight, compatNotes, refQuestions :=
		GenerateAllNarratives(user.Nama, narsisme, machiavellian, psikopati)

	return &models.QuizResult{
		Nama:                user.Nama,
		MBTI:                user.MBTITipe,
		SkorEI:              user.SkorEI,
		SkorSN:              user.SkorSN,
		SkorTF:              user.SkorTF,
		SkorJP:              user.SkorJP,
		Scores:              scores,
		CognitiveStack:      DeriveCognitiveStack(user.MBTITipe),
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
