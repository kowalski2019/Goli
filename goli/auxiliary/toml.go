package auxiliary

import (
	"os"
	"runtime"
	"strings"

	"github.com/laurent22/toml-go"
)

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	osType := runtime.GOOS
	if osType == "windows" || os.Getenv("OS") == "Windows_NT" {
		return "C:\\goli\\config\\config.toml"
	}
	return "/goli/config/config.toml"
}

/*
Get required field from config file.
*/
func GetFromConfig(configField string) string {
	config_path := GetConfigPath()
	var parser toml.Parser

	config := parser.ParseFile(config_path)
	return config.GetString(configField)
}

// getSetupCompleteString properly reads setup_complete as either boolean or string
// This function reads the config file directly to handle boolean values correctly
func getSetupCompleteString() string {
	configPath := GetConfigPath()
	content, err := os.ReadFile(configPath)
	if err != nil {
		return "false"
	}

	lines := strings.Split(string(content), "\n")
	inConstants := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "[constants]" {
			inConstants = true
			continue
		}
		// Check if we're leaving the constants section
		if inConstants && strings.HasPrefix(trimmed, "[") && trimmed != "[constants]" {
			break
		}
		// Look for setup_complete in constants section
		if inConstants && strings.HasPrefix(trimmed, "setup_complete") {
			// Extract the value after =
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				value = strings.Trim(value, `"`)
				value = strings.Trim(value, `'`)
				lower := strings.ToLower(value)
				// Handle boolean values: true, false, 1, 0, yes, no
				if lower == "true" || lower == "1" || lower == "yes" {
					return "true"
				}
				if lower == "false" || lower == "0" || lower == "no" {
					return "false"
				}
				// Return the value as-is if it's something else
				return value
			}
		}
	}

	return "false" // Default to false if not found
}

// GetAllConfig returns all config values as a map
func GetAllConfig() map[string]string {
	config_path := GetConfigPath()
	var parser toml.Parser

	config := parser.ParseFile(config_path)
	result := make(map[string]string)

	// Get known config fields
	result["host"] = config.GetString("constants.host")
	result["port"] = config.GetString("constants.port")
	result["auth_key"] = config.GetString("constants.auth_key")
	result["setup_complete"] = getSetupCompleteString()
	result["setup_password"] = config.GetString("constants.setup_password")
	result["gh_username"] = config.GetString("constants.gh_username")
	result["gh_access_token"] = config.GetString("constants.gh_access_token")
	result["smtp_host"] = config.GetString("constants.smtp_host")
	result["smtp_port"] = config.GetString("constants.smtp_port")
	result["smtp_user"] = config.GetString("constants.smtp_user")
	result["smtp_pass"] = config.GetString("constants.smtp_pass")
	result["smtp_from"] = config.GetString("constants.smtp_from")
	result["smtp_from_name"] = config.GetString("constants.smtp_from_name")

	return result
}

// getCurrentConfigValue extracts the current value of a config key from a line
func getCurrentConfigValue(line string) string {
	parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
	if len(parts) != 2 {
		return ""
	}
	value := strings.TrimSpace(parts[1])
	// Remove quotes if present
	value = strings.Trim(value, `"`)
	value = strings.Trim(value, `'`)
	return value
}

// shouldUpdateField checks if a field should be updated (only if value is different)
func shouldUpdateField(currentValue, newValue string) bool {
	if newValue == "" {
		return false // Don't update if new value is empty
	}
	currentValue = strings.TrimSpace(currentValue)
	newValue = strings.TrimSpace(newValue)
	return currentValue != newValue
}

// UpdateConfig updates the config file with new values, only updating fields that have changed
func UpdateConfig(updates map[string]string) error {
	config_path := GetConfigPath()

	// Get current config to compare values
	currentConfig := GetAllConfig()

	// Filter updates to only include fields that actually changed
	filteredUpdates := make(map[string]string)
	for key, newValue := range updates {
		currentValue := currentConfig[key]
		if shouldUpdateField(currentValue, newValue) {
			filteredUpdates[key] = newValue
		}
	}

	// If no actual changes, return early
	if len(filteredUpdates) == 0 {
		return nil
	}

	// Read existing config file
	content, err := os.ReadFile(config_path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string
	inConstantsSection := false
	hostUpdated := false
	portUpdated := false
	authKeyUpdated := false
	setupCompleteUpdated := false
	setupPasswordUpdated := false
	ghUsernameUpdated := false
	ghAccessTokenUpdated := false
	smtpHostUpdated := false
	smtpPortUpdated := false
	smtpUserUpdated := false
	smtpPassUpdated := false
	smtpFromUpdated := false
	smtpFromNameUpdated := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if we're entering the constants section
		if trimmed == "[constants]" {
			inConstantsSection = true
			updatedLines = append(updatedLines, line)
			continue
		}

		// Check if we're leaving the constants section (new section starts)
		if inConstantsSection && strings.HasPrefix(trimmed, "[") && trimmed != "[constants]" {
			// Add any missing fields before leaving constants section
			if !hostUpdated && filteredUpdates["host"] != "" {
				updatedLines = append(updatedLines, `host = "`+filteredUpdates["host"]+`"`)
				hostUpdated = true
			}
			if !portUpdated && filteredUpdates["port"] != "" {
				updatedLines = append(updatedLines, `port = "`+filteredUpdates["port"]+`"`)
				portUpdated = true
			}
			if !authKeyUpdated && filteredUpdates["auth_key"] != "" {
				updatedLines = append(updatedLines, `auth_key = "`+filteredUpdates["auth_key"]+`"`)
				authKeyUpdated = true
			}
			if !setupCompleteUpdated && filteredUpdates["setup_complete"] != "" {
				updatedLines = append(updatedLines, `setup_complete = `+filteredUpdates["setup_complete"])
				setupCompleteUpdated = true
			}
			if !setupPasswordUpdated && filteredUpdates["setup_password"] != "" {
				updatedLines = append(updatedLines, `setup_password = "`+filteredUpdates["setup_password"]+`"`)
				setupPasswordUpdated = true
			}
			if !ghUsernameUpdated && filteredUpdates["gh_username"] != "" {
				updatedLines = append(updatedLines, `gh_username = "`+filteredUpdates["gh_username"]+`"`)
				ghUsernameUpdated = true
			}
			if !ghAccessTokenUpdated && filteredUpdates["gh_access_token"] != "" {
				updatedLines = append(updatedLines, `gh_access_token = "`+filteredUpdates["gh_access_token"]+`"`)
				ghAccessTokenUpdated = true
			}
			if !smtpHostUpdated && filteredUpdates["smtp_host"] != "" {
				updatedLines = append(updatedLines, `smtp_host = "`+filteredUpdates["smtp_host"]+`"`)
				smtpHostUpdated = true
			}
			if !smtpPortUpdated && filteredUpdates["smtp_port"] != "" {
				updatedLines = append(updatedLines, `smtp_port = "`+filteredUpdates["smtp_port"]+`"`)
				smtpPortUpdated = true
			}
			if !smtpUserUpdated && filteredUpdates["smtp_user"] != "" {
				updatedLines = append(updatedLines, `smtp_user = "`+filteredUpdates["smtp_user"]+`"`)
				smtpUserUpdated = true
			}
			if !smtpPassUpdated && filteredUpdates["smtp_pass"] != "" {
				updatedLines = append(updatedLines, `smtp_pass = "`+filteredUpdates["smtp_pass"]+`"`)
				smtpPassUpdated = true
			}
			if !smtpFromUpdated && filteredUpdates["smtp_from"] != "" {
				updatedLines = append(updatedLines, `smtp_from = "`+filteredUpdates["smtp_from"]+`"`)
				smtpFromUpdated = true
			}
			if !smtpFromNameUpdated && filteredUpdates["smtp_from_name"] != "" {
				updatedLines = append(updatedLines, `smtp_from_name = "`+filteredUpdates["smtp_from_name"]+`"`)
				smtpFromNameUpdated = true
			}
			inConstantsSection = false
		}

		// Update existing fields in constants section
		if inConstantsSection {
			if strings.HasPrefix(trimmed, "host") && !strings.HasPrefix(trimmed, "hostname") && filteredUpdates["host"] != "" {
				updatedLines = append(updatedLines, `host = "`+filteredUpdates["host"]+`"`)
				hostUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "port") && filteredUpdates["port"] != "" {
				updatedLines = append(updatedLines, `port = "`+filteredUpdates["port"]+`"`)
				portUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "auth_key") && filteredUpdates["auth_key"] != "" {
				updatedLines = append(updatedLines, `auth_key = "`+filteredUpdates["auth_key"]+`"`)
				authKeyUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "setup_complete") && filteredUpdates["setup_complete"] != "" {
				updatedLines = append(updatedLines, `setup_complete = `+filteredUpdates["setup_complete"])
				setupCompleteUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "setup_password") {
				if filteredUpdates["setup_password"] != "" {
					updatedLines = append(updatedLines, `setup_password = "`+filteredUpdates["setup_password"]+`"`)
				}
				// If setup_password is empty string, we skip the line (clearing it)
				setupPasswordUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "gh_username") && filteredUpdates["gh_username"] != "" {
				updatedLines = append(updatedLines, `gh_username = "`+filteredUpdates["gh_username"]+`"`)
				ghUsernameUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "gh_access_token") && filteredUpdates["gh_access_token"] != "" {
				updatedLines = append(updatedLines, `gh_access_token = "`+filteredUpdates["gh_access_token"]+`"`)
				ghAccessTokenUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_host") && filteredUpdates["smtp_host"] != "" {
				updatedLines = append(updatedLines, `smtp_host = "`+filteredUpdates["smtp_host"]+`"`)
				smtpHostUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_port") && filteredUpdates["smtp_port"] != "" {
				updatedLines = append(updatedLines, `smtp_port = "`+filteredUpdates["smtp_port"]+`"`)
				smtpPortUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_user") && filteredUpdates["smtp_user"] != "" {
				updatedLines = append(updatedLines, `smtp_user = "`+filteredUpdates["smtp_user"]+`"`)
				smtpUserUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_pass") && filteredUpdates["smtp_pass"] != "" {
				updatedLines = append(updatedLines, `smtp_pass = "`+filteredUpdates["smtp_pass"]+`"`)
				smtpPassUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_from") && !strings.HasPrefix(trimmed, "smtp_from_name") && filteredUpdates["smtp_from"] != "" {
				updatedLines = append(updatedLines, `smtp_from = "`+filteredUpdates["smtp_from"]+`"`)
				smtpFromUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_from_name") && filteredUpdates["smtp_from_name"] != "" {
				updatedLines = append(updatedLines, `smtp_from_name = "`+filteredUpdates["smtp_from_name"]+`"`)
				smtpFromNameUpdated = true
				continue
			}
		}

		updatedLines = append(updatedLines, line)
	}

	// If still in constants section at the end, add missing fields
	if inConstantsSection {
		if !hostUpdated && filteredUpdates["host"] != "" {
			updatedLines = append(updatedLines, `host = "`+filteredUpdates["host"]+`"`)
		}
		if !portUpdated && filteredUpdates["port"] != "" {
			updatedLines = append(updatedLines, `port = "`+filteredUpdates["port"]+`"`)
		}
		if !authKeyUpdated && filteredUpdates["auth_key"] != "" {
			updatedLines = append(updatedLines, `auth_key = "`+filteredUpdates["auth_key"]+`"`)
		}
		if !setupCompleteUpdated && filteredUpdates["setup_complete"] != "" {
			updatedLines = append(updatedLines, `setup_complete = `+filteredUpdates["setup_complete"])
		}
		if !setupPasswordUpdated && filteredUpdates["setup_password"] != "" {
			updatedLines = append(updatedLines, `setup_password = "`+filteredUpdates["setup_password"]+`"`)
		}
		if !ghUsernameUpdated && filteredUpdates["gh_username"] != "" {
			updatedLines = append(updatedLines, `gh_username = "`+filteredUpdates["gh_username"]+`"`)
		}
		if !ghAccessTokenUpdated && filteredUpdates["gh_access_token"] != "" {
			updatedLines = append(updatedLines, `gh_access_token = "`+filteredUpdates["gh_access_token"]+`"`)
		}
		if !smtpHostUpdated && filteredUpdates["smtp_host"] != "" {
			updatedLines = append(updatedLines, `smtp_host = "`+filteredUpdates["smtp_host"]+`"`)
		}
		if !smtpPortUpdated && filteredUpdates["smtp_port"] != "" {
			updatedLines = append(updatedLines, `smtp_port = "`+filteredUpdates["smtp_port"]+`"`)
		}
		if !smtpUserUpdated && filteredUpdates["smtp_user"] != "" {
			updatedLines = append(updatedLines, `smtp_user = "`+filteredUpdates["smtp_user"]+`"`)
		}
		if !smtpPassUpdated && filteredUpdates["smtp_pass"] != "" {
			updatedLines = append(updatedLines, `smtp_pass = "`+filteredUpdates["smtp_pass"]+`"`)
		}
		if !smtpFromUpdated && filteredUpdates["smtp_from"] != "" {
			updatedLines = append(updatedLines, `smtp_from = "`+filteredUpdates["smtp_from"]+`"`)
		}
		if !smtpFromNameUpdated && filteredUpdates["smtp_from_name"] != "" {
			updatedLines = append(updatedLines, `smtp_from_name = "`+filteredUpdates["smtp_from_name"]+`"`)
		}
	} else if !inConstantsSection && len(filteredUpdates) > 0 {
		// No constants section found, add it at the end
		updatedLines = append(updatedLines, "")
		updatedLines = append(updatedLines, "[constants]")
		if filteredUpdates["host"] != "" {
			updatedLines = append(updatedLines, `host = "`+filteredUpdates["host"]+`"`)
		}
		if filteredUpdates["port"] != "" {
			updatedLines = append(updatedLines, `port = "`+filteredUpdates["port"]+`"`)
		}
		if filteredUpdates["auth_key"] != "" {
			updatedLines = append(updatedLines, `auth_key = "`+filteredUpdates["auth_key"]+`"`)
		}
		if filteredUpdates["setup_complete"] != "" {
			updatedLines = append(updatedLines, `setup_complete = `+filteredUpdates["setup_complete"])
		}
		if filteredUpdates["gh_username"] != "" {
			updatedLines = append(updatedLines, `gh_username = "`+filteredUpdates["gh_username"]+`"`)
		}
		if filteredUpdates["gh_access_token"] != "" {
			updatedLines = append(updatedLines, `gh_access_token = "`+filteredUpdates["gh_access_token"]+`"`)
		}
		if filteredUpdates["smtp_host"] != "" {
			updatedLines = append(updatedLines, `smtp_host = "`+filteredUpdates["smtp_host"]+`"`)
		}
		if filteredUpdates["smtp_port"] != "" {
			updatedLines = append(updatedLines, `smtp_port = "`+filteredUpdates["smtp_port"]+`"`)
		}
		if filteredUpdates["smtp_user"] != "" {
			updatedLines = append(updatedLines, `smtp_user = "`+filteredUpdates["smtp_user"]+`"`)
		}
		if filteredUpdates["smtp_pass"] != "" {
			updatedLines = append(updatedLines, `smtp_pass = "`+filteredUpdates["smtp_pass"]+`"`)
		}
		if filteredUpdates["smtp_from"] != "" {
			updatedLines = append(updatedLines, `smtp_from = "`+filteredUpdates["smtp_from"]+`"`)
		}
		if filteredUpdates["smtp_from_name"] != "" {
			updatedLines = append(updatedLines, `smtp_from_name = "`+filteredUpdates["smtp_from_name"]+`"`)
		}
	}

	// Write back to file
	newContent := strings.Join(updatedLines, "\n")
	return os.WriteFile(config_path, []byte(newContent), 0644)
}
