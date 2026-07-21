package repositories

import (
	"ego/database"
	"ego/models"
)

// CreateSession membuat sesi tes baru dan mengembalikan ID
func CreateSession(userID, sessionToken, deviceType, ipAddress string) (string, error) {
	var id string
	query := `INSERT INTO test_sessions (user_id, session_token, device_type, ip_address)
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := database.DB.QueryRow(query, userID, sessionToken, deviceType, ipAddress).Scan(&id)
	return id, err
}

// UpdateSessionCompleted menandai sesi sebagai selesai
func UpdateSessionCompleted(id string) error {
	query := `UPDATE test_sessions SET completed_at = NOW(), is_completed = TRUE WHERE id = $1`
	_, err := database.DB.Exec(query, id)
	return err
}

// GetSessionByID mengambil sesi berdasarkan ID
func GetSessionByID(id string) (*models.Session, error) {
	s := &models.Session{}
	var completedAt *interface{}
	query := `SELECT id, user_id, session_token, started_at, completed_at, device_type, ip_address, is_completed FROM test_sessions WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(
		&s.ID, &s.UserID, &s.SessionToken, &s.StartedAt, &s.CompletedAt,
		&s.DeviceType, &s.IPAddress, &s.IsCompleted,
	)
	_ = completedAt
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetSessionsByUserID mengambil semua sesi untuk seorang user
func GetSessionsByUserID(userID string) ([]models.Session, error) {
	query := `SELECT id, user_id, session_token, started_at, completed_at, device_type, ip_address, is_completed
              FROM test_sessions WHERE user_id = $1 ORDER BY started_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var s models.Session
		err := rows.Scan(
			&s.ID, &s.UserID, &s.SessionToken, &s.StartedAt, &s.CompletedAt,
			&s.DeviceType, &s.IPAddress, &s.IsCompleted,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
