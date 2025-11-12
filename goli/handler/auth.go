package handler

import (
	"crypto/rand"
	"encoding/hex"
	"goli/database"
	"goli/models"
	response_util "goli/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func generateRandomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// LoginHandler verifies username/password and either issues a session or requires 2FA
func LoginHandler(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Channel  string `json:"channel,omitempty"` // optional hint for preferred 2FA channel
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}
	if body.Username == "" || body.Password == "" {
		response_util.SendBadRequestResponseGin(c, "Username and password are required")
		return
	}

	// Get user projection for password and 2FA checks
	userProj, err := getUserProjectionByUsername(body.Username)
	if err != nil {
		response_util.SendUnauthorizedResponseGin(c, "Invalid credentials")
		return
	}

	// Verify password
	if bcrypt.CompareHashAndPassword([]byte(userProj.Password), []byte(body.Password)) != nil {
		response_util.SendUnauthorizedResponseGin(c, "Invalid credentials")
		return
	}

	// 2FA check
	if userTwoFAEnabled(userProj) {
		// send code to selected channel(s)
		preferred := strings.ToLower(strings.TrimSpace(body.Channel))
		now := time.Now().UTC()
		expires := now.Add(10 * time.Minute)
		code := generateNumericCode(6)

		sent := false
		// email
		if userProj.twoFAEmailEnabled && (preferred == "" || preferred == "email") && userProj.Email != "" {
			_, _ = database.CreateTwoFactorCode(userProj.ID, "email", code, expires)
			// Send email asynchronously to avoid blocking the response
			go func(email, code string) {
				if err := response_util.SendEmail2FACode(email, code); err != nil {
					// Log error but don't block response
					// The code is already saved in DB, so user can still verify
					_ = err // Silently log or could use proper logging
				}
			}(userProj.Email, code)
			sent = true
		}
		// sms
		if userProj.twoFASmsEnabled && (preferred == "" || preferred == "sms") && userProj.Phone != "" {
			_, _ = database.CreateTwoFactorCode(userProj.ID, "sms", code, expires)
			// Send SMS asynchronously to avoid blocking the response
			go func(phone, code string) {
				if err := response_util.SendSMS2FACode(phone, code); err != nil {
					// Log error but don't block response
					// The code is already saved in DB, so user can still verify
					_ = err // Silently log or could use proper logging
				}
			}(userProj.Phone, code)
			sent = true
		}

		if !sent {
			response_util.SendBadRequestResponseGin(c, "No valid 2FA delivery channel configured")
			return
		}
		// Return immediately - email/SMS sending happens in background
		response_util.SendJsonResponseGin(c, 200, gin.H{
			"two_fa_required": true,
			"channels":        availableChannels(userProj),
		})
		return
	}

	// No 2FA: issue session
	token, err := generateRandomToken(32)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed generating token")
		return
	}
	expires := time.Now().UTC().Add(24 * time.Hour)
	if _, err := database.CreateSession(userProj.ID, token, expires); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed creating session")
		return
	}
	response_util.SendJsonResponseGin(c, 200, gin.H{
		"token":      token,
		"expires_at": expires,
		"user": gin.H{
			"id":       userProj.ID,
			"username": userProj.Username,
			"email":    userProj.Email,
			"role":     userProj.Role,
		},
	})
}

// Verify2FAHandler verifies a submitted 2FA code and issues a session
func Verify2FAHandler(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Channel  string `json:"channel"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}
	if body.Username == "" || body.Channel == "" || body.Code == "" {
		response_util.SendBadRequestResponseGin(c, "Username, channel and code are required")
		return
	}
	userProj, err := getUserProjectionByUsername(body.Username)
	if err != nil {
		response_util.SendUnauthorizedResponseGin(c, "Invalid credentials")
		return
	}
	tfc, err := database.GetValidTwoFactorCode(userProj.ID, strings.ToLower(body.Channel), body.Code)
	if err != nil {
		response_util.SendUnauthorizedResponseGin(c, "Invalid or expired code")
		return
	}
	_ = database.ConsumeTwoFactorCode(tfc.ID)

	token, err := generateRandomToken(32)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed generating token")
		return
	}
	expires := time.Now().UTC().Add(24 * time.Hour)
	if _, err := database.CreateSession(userProj.ID, token, expires); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed creating session")
		return
	}
	response_util.SendJsonResponseGin(c, 200, gin.H{
		"token":      token,
		"expires_at": expires,
		"user": gin.H{
			"id":       userProj.ID,
			"username": userProj.Username,
			"email":    userProj.Email,
			"role":     userProj.Role,
		},
	})
}

// LogoutHandler deletes the session token
func LogoutHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response_util.SendUnauthorizedResponseGin(c, "Missing Authorization header")
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		response_util.SendUnauthorizedResponseGin(c, "Invalid Authorization header")
		return
	}
	scheme := parts[0]
	cred := parts[1]
	if strings.EqualFold(scheme, "Bearer") {
		_ = database.DeleteSession(cred)
		response_util.SendOkResponseGin(c, "Logged out")
		return
	}
	// Also allow legacy key to "logout" noop
	response_util.SendOkResponseGin(c, "Logged out")
}

// Helpers

func userTwoFAEnabled(user *databaseUserProjection) bool {
	return user.twoFAEmailEnabled || user.twoFASmsEnabled
}

func availableChannels(user *databaseUserProjection) []string {
	ch := []string{}
	if user.twoFAEmailEnabled && user.Email != "" {
		ch = append(ch, "email")
	}
	if user.twoFASmsEnabled && user.Phone != "" {
		ch = append(ch, "sms")
	}
	return ch
}

// databaseUserProjection wraps needed fields from models.User without exposing password
type databaseUserProjection struct {
	ID                int64
	Username          string
	Email             string
	Phone             string
	Role              string
	twoFAEmailEnabled bool
	twoFASmsEnabled   bool
	Password          string
}

// augment GetUserByUsername to a projection (temporary adapter)
func getUserProjectionByUsername(username string) (*databaseUserProjection, error) {
	u, err := database.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return &databaseUserProjection{
		ID:                u.ID,
		Username:          u.Username,
		Email:             u.Email,
		Phone:             u.Phone,
		Role:              u.Role,
		twoFAEmailEnabled: uTwoFAEmail(u),
		twoFASmsEnabled:   uTwoFASms(u),
		Password:          u.PasswordHash, // Use PasswordHash from DB
	}, nil
}

// The following helpers expect the extended fields on models.User; if missing, treat as disabled
func uTwoFAEmail(u *models.User) bool { return u.TwoFAEmailEnabled == 1 }
func uTwoFASms(u *models.User) bool   { return u.TwoFASmsEnabled == 1 }

func generateNumericCode(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := 0; i < n; i++ {
		b[i] = digits[int(b[i])%10]
	}
	return string(b)
}
