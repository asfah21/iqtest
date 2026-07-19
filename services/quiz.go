package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"

	"ego/models"
	"ego/repositories"
)

// ──────────────────────────────────────────────────────────────
// Question Bank — dimuat dari assets/data/questions.json saat startup
// ──────────────────────────────────────────────────────────────

var questionBank []models.QuestionDef

// questionJSON mirrors the JSON structure from questions.json
type questionJSON struct {
	Version   string `json:"version"`
	Questions []struct {
		ID            string            `json:"id"`
		QuestionCode  string            `json:"question_code"`
		Domain        string            `json:"domain"`
		Difficulty    string            `json:"difficulty"`
		Weight        float64           `json:"weight"`
		CorrectOption string            `json:"correct_option"`
		ImagePath     string            `json:"image_path"`
		Options       map[string]string `json:"options"`
	} `json:"questions"`
}

func init() {
	data, err := os.ReadFile("assets/data/questions.json")
	if err != nil {
		panic("gagal membaca assets/data/questions.json: " + err.Error())
	}
	var qj questionJSON
	if err := json.Unmarshal(data, &qj); err != nil {
		panic("gagal parse assets/data/questions.json: " + err.Error())
	}

	for _, q := range qj.Questions {
		def := models.QuestionDef{
			QuestionCode:  q.QuestionCode,
			Domain:        q.Domain,
			Difficulty:    q.Difficulty,
			Weight:        q.Weight,
			CorrectOption: q.CorrectOption,
			ImageURL:      q.ImagePath,
		}
		// Map options in order A, B, C, D
		for i, letter := range []string{"A", "B", "C", "D"} {
			if path, ok := q.Options[letter]; ok {
				def.OptionImages[i] = path
			}
		}
		questionBank = append(questionBank, def)
	}
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
