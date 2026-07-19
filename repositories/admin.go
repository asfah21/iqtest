package repositories

import (
	"ego/database"
	"ego/models"
)

// GetAllUsers mengambil semua data user dari database
func GetAllUsers() ([]models.User, error) {
	query := `SELECT id, nama, email, skor_ei, skor_sn, skor_tf, skor_jp, mbti_tipe, status_pembayaran 
              FROM users_test ORDER BY id DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Nama, &u.Email, &u.SkorEI, &u.SkorSN, &u.SkorTF, &u.SkorJP, &u.MBTITipe, &u.StatusPembayaran)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// GetUserByID mengambil data user berdasarkan ID
func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, nama, email, skor_ei, skor_sn, skor_tf, skor_jp, mbti_tipe, status_pembayaran 
              FROM users_test WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Nama, &user.Email, &user.SkorEI, &user.SkorSN, &user.SkorTF, &user.SkorJP, &user.MBTITipe, &user.StatusPembayaran,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
