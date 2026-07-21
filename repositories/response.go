package repositories

import (
	"ego/database"
	"ego/models"
)

// InsertResponse menyimpan jawaban user untuk satu soal
func InsertResponse(sessionID, questionID string, selectedOption *string, isCorrect bool, timeTakenMs int, timedOut bool) (string, error) {
	var id string
	query := `INSERT INTO session_responses (session_id, question_id, selected_option, is_correct, time_taken_ms, timed_out)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := database.DB.QueryRow(query, sessionID, questionID, selectedOption, isCorrect, timeTakenMs, timedOut).Scan(&id)
	return id, err
}

// InsertResponsesBatch menyimpan semua jawaban dalam satu transaksi
func InsertResponsesBatch(responses []models.SessionResponse, sessionID string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO session_responses (session_id, question_id, selected_option, is_correct, time_taken_ms, timed_out)
                             VALUES ($1, (SELECT id FROM questions WHERE question_code = $2), $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range responses {
		_, err = stmt.Exec(sessionID, r.QuestionCode, r.SelectedOption, r.IsCorrect, r.TimeTakenMs, r.TimedOut)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetSessionResponses mengambil semua jawaban untuk suatu sesi, termasuk info soal
func GetSessionResponses(sessionID string) ([]models.SessionResponse, error) {
	query := `SELECT sr.question_id, q.question_code, sr.selected_option, sr.is_correct, sr.time_taken_ms, sr.timed_out
              FROM session_responses sr
              JOIN questions q ON q.id = sr.question_id
              WHERE sr.session_id = $1
              ORDER BY q.question_code`

	rows, err := database.DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responses []models.SessionResponse
	for rows.Next() {
		var r models.SessionResponse
		err := rows.Scan(&r.QuestionID, &r.QuestionCode, &r.SelectedOption, &r.IsCorrect, &r.TimeTakenMs, &r.TimedOut)
		if err != nil {
			return nil, err
		}
		responses = append(responses, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return responses, nil
}
