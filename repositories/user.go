package repositories

import (
	"ego/database"
	"ego/models"
)

// InsertUser menyimpan data user baru dan mengembalikan ID
func InsertUser(email, nama string, skorEI, skorSN, skorTF, skorJP int, mbtiTipe string) (string, error) {
	var userID string
	query := `INSERT INTO users_test (nama, email, skor_ei, skor_sn, skor_tf, skor_jp, mbti_tipe) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := database.DB.QueryRow(query, nama, email, skorEI, skorSN, skorTF, skorJP, mbtiTipe).Scan(&userID)
	return userID, err
}

// GetUserName mengambil nama user berdasarkan ID
func GetUserName(id string) (string, error) {
	var nama string
	err := database.DB.QueryRow("SELECT nama FROM users_test WHERE id = $1", id).Scan(&nama)
	return nama, err
}

// GetUserResult mengambil data hasil lengkap user (tanpa email)
func GetUserResult(id string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT nama, email, skor_ei, skor_sn, skor_tf, skor_jp, mbti_tipe, status_pembayaran FROM users_test WHERE id = $1"
	err := database.DB.QueryRow(query, id).Scan(
		&user.Nama, &user.Email, &user.SkorEI, &user.SkorSN, &user.SkorTF, &user.SkorJP, &user.MBTITipe, &user.StatusPembayaran,
	)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

// UpdatePaymentStatus memperbarui status pembayaran user menjadi PAID
func UpdatePaymentStatus(id string) error {
	query := "UPDATE users_test SET status_pembayaran = 'PAID' WHERE id = $1"
	_, err := database.DB.Exec(query, id)
	return err
}
