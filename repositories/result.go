package repositories

import (
	"ego/database"
	"ego/models"
)

// InsertIQResult menyimpan hasil IQ test ke database
func InsertIQResult(sessionID string, result models.IQTestResult) (string, error) {
	var id string
	query := `INSERT INTO iq_results (
                session_id, raw_score, max_possible_score,
                mtx_score_pct, seq_score_pct, spa_score_pct, anl_score_pct,
                percentile, estimated_iq, avg_response_ms,
                is_reliable, reliability_flags
              ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
              RETURNING id`

	mtxPct := nullFloat(result.DomainScores["MTX"].Percentage)
	seqPct := nullFloat(result.DomainScores["SEQ"].Percentage)
	spaPct := nullFloat(result.DomainScores["SPA"].Percentage)
	anlPct := nullFloat(result.DomainScores["ANL"].Percentage)

	err := database.DB.QueryRow(query,
		sessionID, result.RawScore, result.MaxPossible,
		mtxPct, seqPct, spaPct, anlPct,
		result.Percentile, result.EstimatedIQ, result.AvgResponseMs,
		result.IsReliable, result.ReliabilityFlags,
	).Scan(&id)
	return id, err
}

// GetIQResultBySession mengambil hasil IQ berdasarkan session_id
func GetIQResultBySession(sessionID string) (*models.IQTestResult, error) {
	r := &models.IQTestResult{}
	var mtxPct, seqPct, spaPct, anlPct *float64
	query := `SELECT raw_score, max_possible_score,
              mtx_score_pct, seq_score_pct, spa_score_pct, anl_score_pct,
              percentile, estimated_iq, avg_response_ms,
              is_reliable
              FROM iq_results WHERE session_id = $1`
	err := database.DB.QueryRow(query, sessionID).Scan(
		&r.RawScore, &r.MaxPossible,
		&mtxPct, &seqPct, &spaPct, &anlPct,
		&r.Percentile, &r.EstimatedIQ, &r.AvgResponseMs,
		&r.IsReliable,
	)
	if err != nil {
		return nil, err
	}

	r.DomainScores = map[string]models.DomainScore{
		"MTX": {Domain: "MTX", Percentage: safeDeref(mtxPct)},
		"SEQ": {Domain: "SEQ", Percentage: safeDeref(seqPct)},
		"SPA": {Domain: "SPA", Percentage: safeDeref(spaPct)},
		"ANL": {Domain: "ANL", Percentage: safeDeref(anlPct)},
	}

	return r, nil
}

// GetAllRawScores mengambil semua raw_score dari iq_results untuk kalkulasi persentil
func GetAllRawScores() ([]float64, error) {
	query := `SELECT raw_score FROM iq_results`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []float64
	for rows.Next() {
		var s float64
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}

	return scores, nil
}

// nullFloat mengembalikan *float64 atau nil jika value == 0
func nullFloat(v float64) *float64 {
	if v == 0 {
		return nil
	}
	return &v
}

// safeDeref mengembalikan nilai float64 dari pointer, atau 0 jika nil
func safeDeref(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
