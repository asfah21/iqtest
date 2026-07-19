package repositories

import (
	"ego/database"
	"ego/models"
)

// InsertPayment menyimpan record pembayaran baru
func InsertPayment(userID, sessionID string, amount float64, currency, status string) (string, error) {
	var id string
	query := `INSERT INTO payments (user_id, session_id, amount, currency, status)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := database.DB.QueryRow(query, userID, sessionID, amount, currency, status).Scan(&id)
	return id, err
}

// UpdatePaymentStatus memperbarui status pembayaran menjadi PAID
func UpdatePaymentStatus(sessionID string) error {
	query := `UPDATE payments SET status = 'PAID', paid_at = NOW() WHERE session_id = $1`
	_, err := database.DB.Exec(query, sessionID)
	return err
}

// GetPaymentBySession mengambil data pembayaran berdasarkan session_id
func GetPaymentBySession(sessionID string) (*models.Payment, error) {
	p := &models.Payment{}
	query := `SELECT id, user_id, session_id, amount, currency, status, payment_method, paid_at, created_at
              FROM payments WHERE session_id = $1`
	err := database.DB.QueryRow(query, sessionID).Scan(
		&p.ID, &p.UserID, &p.SessionID, &p.Amount, &p.Currency,
		&p.Status, &p.PaymentMethod, &p.PaidAt, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}
