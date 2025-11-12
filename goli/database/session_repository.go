package database

import (
	"time"
)

type Session struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// CreateSession inserts a new session for a user
func CreateSession(userID int64, token string, expiresAt time.Time) (*Session, error) {
	query := `INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)
			  RETURNING id, created_at`
	s := &Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	err := DB.QueryRow(query, userID, token, expiresAt).Scan(&s.ID, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetSessionByToken retrieves a session by its token
func GetSessionByToken(token string) (*Session, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = ?`
	s := &Session{}
	err := DB.QueryRow(query, token).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteSession deletes a session by token
func DeleteSession(token string) error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	return err
}

// CleanupExpiredSessions removes expired sessions
func CleanupExpiredSessions() error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE expires_at < ?`, time.Now().UTC())
	return err
}
