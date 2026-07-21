package repositories

import (
	"database/sql"

	"ego/database"
	"ego/models"
)

// AdminUserRow menyimpan data user dengan skor dari iq_results untuk dashboard admin
type AdminUserRow struct {
	models.User
	RawScore    sql.NullFloat64
	Percentile  sql.NullFloat64
	EstimatedIQ sql.NullFloat64
}

// GetAllUsers mengambil semua data user dengan LEFT JOIN ke iq_results untuk raw_score, percentile, estimated_iq
func GetAllUsers() ([]AdminUserRow, error) {
	query := `SELECT u.id, u.email, u.nama, u.phone, u.created_at, u.updated_at,
                     r.raw_score, r.percentile, r.estimated_iq
              FROM users u
              LEFT JOIN iq_results r ON r.session_id IN (
                  SELECT id FROM test_sessions WHERE user_id = u.id LIMIT 1
              )
              ORDER BY u.created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []AdminUserRow
	for rows.Next() {
		var u AdminUserRow
		err := rows.Scan(
			&u.ID, &u.Email, &u.Nama, &u.Phone, &u.CreatedAt, &u.UpdatedAt,
			&u.RawScore, &u.Percentile, &u.EstimatedIQ,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
