package repositories

import (
	"ego/database"
	"ego/models"
)

// InsertUser menyimpan data user baru dan mengembalikan ID
func InsertUser(email, nama, phone string) (string, error) {
	var userID string
	query := `INSERT INTO users (email, nama, phone) 
              VALUES ($1, $2, $3) RETURNING id`
	err := database.DB.QueryRow(query, email, nama, phone).Scan(&userID)
	return userID, err
}

// GetUserByID mengambil data user berdasarkan ID
func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, nama, phone, created_at, updated_at FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Nama, &user.Phone, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail mencari user berdasarkan email (untuk handle retake)
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, nama, phone, created_at, updated_at FROM users WHERE email = $1`
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Nama, &user.Phone, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
