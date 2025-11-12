package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	aux "goli/auxiliary"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SendJsonResponse sends a JSON response
func SendJsonResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

// SendOkResponse sends a success response
func SendOkResponse(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"status":      "success",
		"description": message,
	}
	jsonData, _ := json.Marshal(response)
	SendJsonResponse(w, http.StatusOK, jsonData)
}

// SendBadRequestResponse sends a bad request error response
func SendBadRequestResponse(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"status":      "error",
		"description": message,
	}
	jsonData, _ := json.Marshal(response)
	SendJsonResponse(w, http.StatusBadRequest, jsonData)
}

// SendUnauthorizedResponse sends an unauthorized error response
func SendUnauthorizedResponse(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"status":      "error",
		"description": message,
	}
	jsonData, _ := json.Marshal(response)
	SendJsonResponse(w, http.StatusUnauthorized, jsonData)
}

// SendNotFoundResponse sends a not found error response
func SendNotFoundResponse(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"status":      "error",
		"description": message,
	}
	jsonData, _ := json.Marshal(response)
	SendJsonResponse(w, http.StatusNotFound, jsonData)
}

// SendInternalServerErrorResponse sends an internal server error response
func SendInternalServerErrorResponse(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"status":      "error",
		"description": message,
	}
	jsonData, _ := json.Marshal(response)
	SendJsonResponse(w, http.StatusInternalServerError, jsonData)
}

// Gin-compatible response functions

// SendJsonResponseGin sends a JSON response using Gin context
func SendJsonResponseGin(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// SendOkResponseGin sends a success response using Gin context
func SendOkResponseGin(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"description": message,
	})
}

// SendBadRequestResponseGin sends a bad request error response using Gin context
func SendBadRequestResponseGin(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":      "error",
		"description": message,
	})
}

// SendUnauthorizedResponseGin sends an unauthorized error response using Gin context
func SendUnauthorizedResponseGin(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":      "error",
		"description": message,
	})
}

// SendNotFoundResponseGin sends a not found error response using Gin context
func SendNotFoundResponseGin(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":      "error",
		"description": message,
	})
}

// SendInternalServerErrorResponseGin sends an internal server error response using Gin context
func SendInternalServerErrorResponseGin(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":      "error",
		"description": message,
	})
}

// Notification helpers for 2FA

// SendEmail2FACode sends a 2FA code via email using SMTP configuration in config.toml
// This function now uses a timeout and supports TLS/SSL connections
func SendEmail2FACode(toEmail, code string) error {
	smtpHost := aux.GetFromConfig("constants.smtp_host")
	smtpPortStr := aux.GetFromConfig("constants.smtp_port")
	smtpUser := aux.GetFromConfig("constants.smtp_user")
	smtpPass := aux.GetFromConfig("constants.smtp_pass")
	from := aux.GetFromConfig("constants.smtp_from")
	fromName := aux.GetFromConfig("constants.smtp_from_name")
	if fromName == "" {
		fromName = "Goli"
	}

	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPass == "" || from == "" {
		// Fallback: log only
		fmt.Printf("[2FA EMAIL] To: %s Code: %s (SMTP not configured)\n", toEmail, code)
		return nil
	}

	// Parse port
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	// Use context with timeout to prevent blocking indefinitely (30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Channel to receive the result
	resultChan := make(chan error, 1)

	go func() {
		err := sendEmailWithTLS(ctx, smtpHost, smtpPort, smtpUser, smtpPass, from, fromName, toEmail, code)
		resultChan <- err
	}()

	// Wait for result or timeout
	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("SMTP send timeout after 30 seconds")
	}
}

// sendEmailWithTLS sends an email with proper TLS/SSL support
func sendEmailWithTLS(ctx context.Context, host string, port int, username, password, from, fromName, to, code string) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", username, password, host)

	// Build proper email message with headers
	message := buildEmailMessage(from, fromName, to, "Your Goli verification code", code)

	// Determine if we should use SSL (port 465) or STARTTLS (port 587, 25, etc.)
	useSSL := port == 465
	useTLS := port == 587 || port == 25 || port == 2525

	// Use Dialer with timeout for connection
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	var conn net.Conn
	var err error

	if useSSL {
		// Direct SSL/TLS connection (port 465)
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		}
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server via SSL: %w", err)
		}
	} else {
		// Plain TCP connection (will upgrade to TLS if needed)
		conn, err = dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer func() {
		if quitErr := client.Quit(); quitErr != nil {
			// Log but don't fail on quit error
			_ = quitErr
		}
	}()

	// Upgrade to TLS if needed (STARTTLS)
	if useTLS && !useSSL {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email data
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data connection: %w", err)
	}

	_, err = w.Write(message)
	if err != nil {
		w.Close()
		return fmt.Errorf("failed to write email data: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data connection: %w", err)
	}

	return nil
}

// buildEmailMessage builds a properly formatted email message with headers
func buildEmailMessage(from, fromName, to, subject, body string) []byte {
	var message bytes.Buffer

	// Email headers
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	message.WriteString("\r\n")

	// Email body
	message.WriteString(fmt.Sprintf("Your verification code is: %s\r\n\r\n", body))
	message.WriteString("This code will expire in 10 minutes.\r\n")
	message.WriteString("If you did not request this code, please ignore this email.\r\n")

	return message.Bytes()
}

// SendSMS2FACode sends a 2FA code via SMS; if SMS provider not configured, logs to console
// This function now uses a timeout to prevent blocking indefinitely
func SendSMS2FACode(toNumber, code string) error {
	twSid := aux.GetFromConfig("constants.twilio_sid")
	twToken := aux.GetFromConfig("constants.twilio_token")
	from := aux.GetFromConfig("constants.twilio_from")
	if twSid == "" || twToken == "" || from == "" {
		fmt.Printf("[2FA SMS] To: %s Code: %s (SMS not configured)\n", toNumber, code)
		return nil
	}
	// Simple Twilio API call via HTTP with timeout
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + twSid + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", toNumber)
	msgData.Set("From", from)
	msgData.Set("Body", "Your Goli verification code is: "+code)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// HTTP client with 30 second timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		return err
	}
	req.SetBasicAuth(twSid, twToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("twilio send failed with status %s", resp.Status)
}
