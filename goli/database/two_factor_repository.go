package database

import (
	"database/sql"
	"time"
)

type TwoFactorCode struct {
	ID        int64
	UserID    int64
	Channel   string
	Code      string
	ExpiresAt time.Time
	Consumed  bool
	CreatedAt time.Time
}

// CreateTwoFactorCode saves a new 2FA code
func CreateTwoFactorCode(userID int64, channel, code string, expiresAt time.Time) (*TwoFactorCode, error) {
	query := `INSERT INTO two_factor_codes (user_id, channel, code, expires_at) 
			  VALUES (?, ?, ?, ?) RETURNING id, created_at`
	tfc := &TwoFactorCode{
		UserID:    userID,
		Channel:   channel,
		Code:      code,
		ExpiresAt: expiresAt,
	}
	err := DB.QueryRow(query, userID, channel, code, expiresAt).Scan(&tfc.ID, &tfc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tfc, nil
}

// GetValidTwoFactorCode returns a valid, unconsumed code if present
func GetValidTwoFactorCode(userID int64, channel, code string) (*TwoFactorCode, error) {
	query := `SELECT id, user_id, channel, code, expires_at, consumed, created_at
			  FROM two_factor_codes
			  WHERE user_id = ? AND channel = ? AND code = ? AND consumed = 0`
	tfc := &TwoFactorCode{}
	err := DB.QueryRow(query, userID, channel, code).Scan(
		&tfc.ID, &tfc.UserID, &tfc.Channel, &tfc.Code, &tfc.ExpiresAt, &tfc.Consumed, &tfc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if time.Now().After(tfc.ExpiresAt) {
		return nil, sql.ErrNoRows
	}
	return tfc, nil
}

// ConsumeTwoFactorCode marks a code as consumed
func ConsumeTwoFactorCode(id int64) error {
	_, err := DB.Exec(`UPDATE two_factor_codes SET consumed = 1 WHERE id = ?`, id)
	return err
}

// CleanupExpiredTwoFactorCodes removes expired codes
func CleanupExpiredTwoFactorCodes() error {
	_, err := DB.Exec(`DELETE FROM two_factor_codes WHERE expires_at < ? OR consumed = 1`, time.Now().UTC())
	return err
}