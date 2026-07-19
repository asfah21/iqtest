package repositories

import (
	"ego/database"
	"ego/models"
)

// InsertQuestion menyimpan soal baru ke database
func InsertQuestion(q models.QuestionDef) (string, error) {
	var id string
	query := `INSERT INTO questions (question_code, domain, difficulty, weight, image_url,
               option_a_image, option_b_image, option_c_image, option_d_image, correct_option,
               p_value, discrimination, is_active)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
              RETURNING id`
	err := database.DB.QueryRow(query,
		q.QuestionCode, q.Domain, q.Difficulty, q.Weight, q.ImageURL,
		q.OptionImages[0], q.OptionImages[1], q.OptionImages[2], q.OptionImages[3],
		q.CorrectOption, q.PValue, q.Discrimination, true,
	).Scan(&id)
	return id, err
}

// GetActiveQuestions mengambil semua soal aktif yang sudah di database
func GetActiveQuestions() ([]models.QuestionDef, error) {
	query := `SELECT id, question_code, domain, difficulty, weight, image_url,
              option_a_image, option_b_image, option_c_image, option_d_image, correct_option,
              p_value, discrimination
              FROM questions WHERE is_active = TRUE ORDER BY question_code`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.QuestionDef
	for rows.Next() {
		var q models.QuestionDef
		var optA, optB, optC, optD string
		err := rows.Scan(
			&q.ID, &q.QuestionCode, &q.Domain, &q.Difficulty, &q.Weight, &q.ImageURL,
			&optA, &optB, &optC, &optD, &q.CorrectOption,
			&q.PValue, &q.Discrimination,
		)
		if err != nil {
			return nil, err
		}
		q.OptionImages = [4]string{optA, optB, optC, optD}
		questions = append(questions, q)
	}

	return questions, nil
}

// GetQuestionByCode mengambil satu soal berdasarkan question_code
func GetQuestionByCode(code string) (*models.QuestionDef, error) {
	q := &models.QuestionDef{}
	var optA, optB, optC, optD string
	query := `SELECT id, question_code, domain, difficulty, weight, image_url,
              option_a_image, option_b_image, option_c_image, option_d_image, correct_option,
              p_value, discrimination
              FROM questions WHERE question_code = $1`
	err := database.DB.QueryRow(query, code).Scan(
		&q.ID, &q.QuestionCode, &q.Domain, &q.Difficulty, &q.Weight, &q.ImageURL,
		&optA, &optB, &optC, &optD, &q.CorrectOption,
		&q.PValue, &q.Discrimination,
	)
	if err != nil {
		return nil, err
	}
	q.OptionImages = [4]string{optA, optB, optC, optD}
	return q, nil
}
