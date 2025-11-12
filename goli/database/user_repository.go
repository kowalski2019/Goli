package database

import (
	"goli/models"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user
func CreateUser(user *models.User) (*models.User, error) {
	query := `INSERT INTO users (username, email, phone, password, role, two_fa_email_enabled, two_fa_sms_enabled) 
			  VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, created_at, updated_at`

	// Hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = DB.QueryRow(query, user.Username, user.Email, user.Phone, string(hashedPassword), user.Role, user.TwoFAEmailEnabled, user.TwoFASmsEnabled).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func GetUser(id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, phone, password, role, two_fa_email_enabled, two_fa_sms_enabled, created_at, updated_at 
			  FROM users WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.Role, &user.TwoFAEmailEnabled, &user.TwoFASmsEnabled, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ListUsers retrieves all users (without passwords)
func ListUsers() ([]*models.User, error) {
	query := `SELECT id, username, email, phone, role, two_fa_email_enabled, two_fa_sms_enabled, created_at, updated_at 
			  FROM users ORDER BY created_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Phone,
			&user.Role, &user.TwoFAEmailEnabled, &user.TwoFASmsEnabled, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates an existing user
func UpdateUser(user *models.User, updatePassword bool) error {
	var query string
	var err error

	if updatePassword {
		// Hash new password
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}

		query = `UPDATE users SET email = ?, phone = ?, password = ?, role = ?, two_fa_email_enabled = ?, two_fa_sms_enabled = ?, updated_at = CURRENT_TIMESTAMP 
				 WHERE id = ?`
		_, err = DB.Exec(query, user.Email, user.Phone, string(hashedPassword), user.Role, user.TwoFAEmailEnabled, user.TwoFASmsEnabled, user.ID)
	} else {
		query = `UPDATE users SET email = ?, phone = ?, role = ?, two_fa_email_enabled = ?, two_fa_sms_enabled = ?, updated_at = CURRENT_TIMESTAMP 
				 WHERE id = ?`
		_, err = DB.Exec(query, user.Email, user.Phone, user.Role, user.TwoFAEmailEnabled, user.TwoFASmsEnabled, user.ID)
	}
	return err
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, phone, password, role, two_fa_email_enabled, two_fa_sms_enabled, created_at, updated_at 
			  FROM users WHERE username = ?`

	err := DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.Role, &user.TwoFAEmailEnabled, &user.TwoFASmsEnabled, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser deletes a user by ID
func DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}
