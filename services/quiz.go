package services

import (
	"crypto/rand"
	"encoding/hex"
	"math"

	"ego/models"
	"ego/repositories"
)

// ──────────────────────────────────────────────────────────────
// Question Definition — metadata setiap soal IQ Test
// ──────────────────────────────────────────────────────────────

var questionBank = []models.QuestionDef{
	{QuestionCode: "Q_MTX_001", Domain: "MTX", Difficulty: "easy", Weight: 1.0, CorrectOption: "A"},
	{QuestionCode: "Q_MTX_002", Domain: "MTX", Difficulty: "easy", Weight: 1.0, CorrectOption: "B"},
	{QuestionCode: "Q_MTX_003", Domain: "MTX", Difficulty: "medium", Weight: 1.5, CorrectOption: "C"},
	{QuestionCode: "Q_MTX_004", Domain: "MTX", Difficulty: "medium", Weight: 1.5, CorrectOption: "D"},
	{QuestionCode: "Q_MTX_005", Domain: "MTX", Difficulty: "hard", Weight: 2.0, CorrectOption: "A"},
	{QuestionCode: "Q_MTX_006", Domain: "MTX", Difficulty: "very_hard", Weight: 2.5, CorrectOption: "B"},
	{QuestionCode: "Q_SEQ_001", Domain: "SEQ", Difficulty: "easy", Weight: 1.0, CorrectOption: "A"},
	{QuestionCode: "Q_SEQ_002", Domain: "SEQ", Difficulty: "medium", Weight: 1.5, CorrectOption: "B"},
	{QuestionCode: "Q_SEQ_003", Domain: "SEQ", Difficulty: "medium", Weight: 1.5, CorrectOption: "C"},
	{QuestionCode: "Q_SEQ_004", Domain: "SEQ", Difficulty: "hard", Weight: 2.0, CorrectOption: "D"},
	{QuestionCode: "Q_SEQ_005", Domain: "SEQ", Difficulty: "hard", Weight: 2.0, CorrectOption: "A"},
	{QuestionCode: "Q_SPA_001", Domain: "SPA", Difficulty: "medium", Weight: 1.5, CorrectOption: "B"},
	{QuestionCode: "Q_SPA_002", Domain: "SPA", Difficulty: "medium", Weight: 1.5, CorrectOption: "C"},
	{QuestionCode: "Q_SPA_003", Domain: "SPA", Difficulty: "hard", Weight: 2.0, CorrectOption: "D"},
	{QuestionCode: "Q_SPA_004", Domain: "SPA", Difficulty: "very_hard", Weight: 2.5, CorrectOption: "A"},
	{QuestionCode: "Q_SPA_005", Domain: "SPA", Difficulty: "very_hard", Weight: 2.5, CorrectOption: "B"},
	{QuestionCode: "Q_ANL_001", Domain: "ANL", Difficulty: "easy", Weight: 1.0, CorrectOption: "C"},
	{QuestionCode: "Q_ANL_002", Domain: "ANL", Difficulty: "medium", Weight: 1.5, CorrectOption: "D"},
	{QuestionCode: "Q_ANL_003", Domain: "ANL", Difficulty: "medium", Weight: 1.5, CorrectOption: "A"},
	{QuestionCode: "Q_ANL_004", Domain: "ANL", Difficulty: "hard", Weight: 2.0, CorrectOption: "B"},
}

// ──────────────────────────────────────────────────────────────
// CalculateIQResult
// ──────────────────────────────────────────────────────────────

func CalculateIQResult(responses []models.SessionResponse) models.IQTestResult {
	domainScores := map[string]float64{"MTX": 0, "SEQ": 0, "SPA": 0, "ANL": 0}
	domainMax := map[string]float64{"MTX": 0, "SEQ": 0, "SPA": 0, "ANL": 0}
	totalTimeMs := 0

	for _, response := range responses {
		for _, q := range questionBank {
			if q.QuestionCode != response.QuestionCode {
				continue
			}
			domainMax[q.Domain] += q.Weight
			if response.IsCorrect {
				domainScores[q.Domain] += q.Weight
			}
			totalTimeMs += response.TimeTakenMs
			break
		}
	}

	rawScore := domainScores["MTX"] + domainScores["SEQ"] + domainScores["SPA"] + domainScores["ANL"]
	maxPossible := 30.5

	domainPct := make(map[string]models.DomainScore)
	for _, domain := range []string{"MTX", "SEQ", "SPA", "ANL"} {
		pct := 0.0
		if domainMax[domain] > 0 {
			pct = domainScores[domain] / domainMax[domain] * 100
		}
		domainPct[domain] = models.DomainScore{
			Domain:      domain,
			RawScore:    domainScores[domain],
			MaxPossible: domainMax[domain],
			Percentage:  pct,
		}
	}

	return models.IQTestResult{
		RawScore:         rawScore,
		MaxPossible:      maxPossible,
		DomainScores:     domainPct,
		AvgResponseMs:    totalTimeMs / len(responses),
		IsReliable:       true,
		ReliabilityFlags: nil,
	}
}

// ──────────────────────────────────────────────────────────────
// CalculatePercentile
// ──────────────────────────────────────────────────────────────

func CalculatePercentile(rawScore float64, allScores []float64) float64 {
	if len(allScores) == 0 {
		return 0
	}
	count := 0
	for _, s := range allScores {
		if s <= rawScore {
			count++
		}
	}
	return float64(count) / float64(len(allScores)) * 100
}

// ──────────────────────────────────────────────────────────────
// EstimateIQ — Deviation IQ (Wechsler-style)
// ──────────────────────────────────────────────────────────────

func EstimateIQ(rawScore float64, populationMean float64, populationStdDev float64) *float64 {
	if populationStdDev == 0 {
		return nil
	}
	zScore := (rawScore - populationMean) / populationStdDev
	iq := 100 + (zScore * 15)
	iq = math.Round(iq*10) / 10
	return &iq
}

// ──────────────────────────────────────────────────────────────
// Anti-Cheating Detection (Per IQTEST.md §9)
// ──────────────────────────────────────────────────────────────

func DetectSpeedGuessing(responses []models.SessionResponse) bool {
	if len(responses) == 0 {
		return false
	}
	fastCount := 0
	correctCount := 0
	for _, r := range responses {
		if r.TimeTakenMs > 0 && r.TimeTakenMs < 3000 {
			fastCount++
			if r.IsCorrect {
				correctCount++
			}
		}
	}
	if fastCount == 0 {
		return false
	}
	fastRatio := float64(fastCount) / float64(len(responses))
	accuracy := float64(correctCount) / float64(fastCount)
	return fastRatio > 0.3 && accuracy < 0.25
}

func DetectStraightPattern(responses []models.SessionResponse) bool {
	optionCounts := map[string]int{"A": 0, "B": 0, "C": 0, "D": 0}
	for _, r := range responses {
		if r.SelectedOption != nil {
			optionCounts[*r.SelectedOption]++
		}
	}
	for _, count := range optionCounts {
		if count >= 15 {
			return true
		}
	}
	return false
}

func AssessReliability(responses []models.SessionResponse, tabSwitchCount int) (bool, []string) {
	var reasons []string
	if DetectSpeedGuessing(responses) {
		reasons = append(reasons, "speed_guessing")
	}
	if DetectStraightPattern(responses) {
		reasons = append(reasons, "straight_pattern")
	}
	if tabSwitchCount > 3 {
		reasons = append(reasons, "tab_switch_excessive")
	}
	return len(reasons) == 0, reasons
}

// ──────────────────────────────────────────────────────────────
// GetQuestions
// ──────────────────────────────────────────────────────────────

func GetQuestions() []models.QuestionDef {
	return questionBank
}

func generateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ──────────────────────────────────────────────────────────────
// ProcessQuizAnswers
// ──────────────────────────────────────────────────────────────

func ProcessQuizAnswers(email, nama string, rawAnswers map[string]string, tabSwitchCount int) (string, error) {
	var responses []models.SessionResponse
	for _, q := range questionBank {
		selectedOption, exists := rawAnswers[q.QuestionCode]
		isCorrect := exists && selectedOption == q.CorrectOption
		resp := models.SessionResponse{
			QuestionCode: q.QuestionCode,
			IsCorrect:    isCorrect,
		}
		if exists {
			resp.SelectedOption = &selectedOption
		}
		responses = append(responses, resp)
	}

	result := CalculateIQResult(responses)

	allScores, err := repositories.GetAllRawScores()
	if err != nil {
		allScores = nil
	}
	percentile := CalculatePercentile(result.RawScore, allScores)
	result.Percentile = &percentile
	result.EstimatedIQ = EstimateIQ(result.RawScore, 0, 0)

	// Assess reliability
	isReliable, reasons := AssessReliability(responses, tabSwitchCount)
	result.IsReliable = isReliable
	result.ReliabilityFlags = reasons

	var userID string
	existingUser, err := repositories.GetUserByEmail(email)
	if err != nil {
		userID, err = repositories.InsertUser(email, nama, "")
		if err != nil {
			return "", err
		}
	} else {
		userID = existingUser.ID
	}

	sessionToken := generateSessionToken()
	sessionID, err := repositories.CreateSession(userID, sessionToken, "web", "")
	if err != nil {
		return "", err
	}

	err = repositories.InsertResponsesBatch(responses, sessionID)
	if err != nil {
		return "", err
	}

	_, err = repositories.InsertIQResult(sessionID, result)
	if err != nil {
		return "", err
	}

	repositories.UpdateSessionCompleted(sessionID)
	return sessionID, nil
}

// ──────────────────────────────────────────────────────────────
// GetPaywallData
// ──────────────────────────────────────────────────────────────

func GetPaywallData(sessionID string) (*models.PaywallData, error) {
	session, err := repositories.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}
	user, err := repositories.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}
	return &models.PaywallData{ID: sessionID, Nama: user.Nama}, nil
}

// ──────────────────────────────────────────────────────────────
// GetQuizResult
// ──────────────────────────────────────────────────────────────

func GetQuizResult(sessionID string) (*models.QuizResult, error) {
	payment, err := repositories.GetPaymentBySession(sessionID)
	if err != nil || payment.Status != "PAID" {
		return nil, nil
	}

	result, err := repositories.GetIQResultBySession(sessionID)
	if err != nil {
		return nil, nil
	}

	session, err := repositories.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}

	user, err := repositories.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}

	percentile := 0.0
	if result.Percentile != nil {
		percentile = *result.Percentile
	}
	execSummary, kekuatan, areaPerhatian, _ := GenerateAllNarratives(
		user.Nama, result.RawScore, result.MaxPossible, result.DomainScores, percentile, result.EstimatedIQ,
	)

	return &models.QuizResult{
		ID:               sessionID,
		Nama:             user.Nama,
		RawScore:         result.RawScore,
		MaxPossible:      result.MaxPossible,
		Percentile:       result.Percentile,
		EstimatedIQ:      result.EstimatedIQ,
		DomainScores:     result.DomainScores,
		AvgResponseMs:    result.AvgResponseMs,
		IsReliable:       result.IsReliable,
		ReliabilityFlags: result.ReliabilityFlags,
		ExecutiveSummary: execSummary,
		Kekuatan:         kekuatan,
		AreaPerhatian:    areaPerhatian,
	}, nil
}

// ──────────────────────────────────────────────────────────────
// ConfirmPayment
// ──────────────────────────────────────────────────────────────

func ConfirmPayment(sessionID string) error {
	return repositories.UpdatePaymentStatus(sessionID)
}

// ──────────────────────────────────────────────────────────────
// GetAllUsers
// ──────────────────────────────────────────────────────────────

func GetAllUsers() ([]repositories.AdminUserRow, error) {
	return repositories.GetAllUsers()
}

// ──────────────────────────────────────────────────────────────
// GetUserByID
// ──────────────────────────────────────────────────────────────

func GetUserByID(id string) (*models.User, error) {
	return repositories.GetUserByID(id)
}
