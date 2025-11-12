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

// GetAllConfig returns all config values as a map
func GetAllConfig() map[string]string {
	config_path := GetConfigPath()
	var parser toml.Parser

	config := parser.ParseFile(config_path)
	result := make(map[string]string)

	// Get known config fields
	result["port"] = config.GetString("constants.port")
	result["auth_key"] = config.GetString("constants.auth_key")
	result["setup_complete"] = config.GetString("constants.setup_complete")
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

// UpdateConfig updates the config file with new values
func UpdateConfig(updates map[string]string) error {
	config_path := GetConfigPath()

	// Read existing config file
	content, err := os.ReadFile(config_path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string
	inConstantsSection := false
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
			if !portUpdated && updates["port"] != "" {
				updatedLines = append(updatedLines, `port = "`+updates["port"]+`"`)
				portUpdated = true
			}
			if !authKeyUpdated && updates["auth_key"] != "" {
				updatedLines = append(updatedLines, `auth_key = "`+updates["auth_key"]+`"`)
				authKeyUpdated = true
			}
			if !setupCompleteUpdated && updates["setup_complete"] != "" {
				updatedLines = append(updatedLines, `setup_complete = `+updates["setup_complete"])
				setupCompleteUpdated = true
			}
			if !setupPasswordUpdated && updates["setup_password"] != "" {
				updatedLines = append(updatedLines, `setup_password = "`+updates["setup_password"]+`"`)
				setupPasswordUpdated = true
			}
			if !ghUsernameUpdated && updates["gh_username"] != "" {
				updatedLines = append(updatedLines, `gh_username = "`+updates["gh_username"]+`"`)
				ghUsernameUpdated = true
			}
			if !ghAccessTokenUpdated && updates["gh_access_token"] != "" {
				updatedLines = append(updatedLines, `gh_access_token = "`+updates["gh_access_token"]+`"`)
				ghAccessTokenUpdated = true
			}
			if !smtpHostUpdated && updates["smtp_host"] != "" {
				updatedLines = append(updatedLines, `smtp_host = "`+updates["smtp_host"]+`"`)
				smtpHostUpdated = true
			}
			if !smtpPortUpdated && updates["smtp_port"] != "" {
				updatedLines = append(updatedLines, `smtp_port = "`+updates["smtp_port"]+`"`)
				smtpPortUpdated = true
			}
			if !smtpUserUpdated && updates["smtp_user"] != "" {
				updatedLines = append(updatedLines, `smtp_user = "`+updates["smtp_user"]+`"`)
				smtpUserUpdated = true
			}
			if !smtpPassUpdated && updates["smtp_pass"] != "" {
				updatedLines = append(updatedLines, `smtp_pass = "`+updates["smtp_pass"]+`"`)
				smtpPassUpdated = true
			}
			if !smtpFromUpdated && updates["smtp_from"] != "" {
				updatedLines = append(updatedLines, `smtp_from = "`+updates["smtp_from"]+`"`)
				smtpFromUpdated = true
			}
			if !smtpFromNameUpdated && updates["smtp_from_name"] != "" {
				updatedLines = append(updatedLines, `smtp_from_name = "`+updates["smtp_from_name"]+`"`)
				smtpFromNameUpdated = true
			}
			inConstantsSection = false
		}

		// Update existing fields in constants section
		if inConstantsSection {
			if strings.HasPrefix(trimmed, "port") && updates["port"] != "" {
				updatedLines = append(updatedLines, `port = "`+updates["port"]+`"`)
				portUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "auth_key") && updates["auth_key"] != "" {
				updatedLines = append(updatedLines, `auth_key = "`+updates["auth_key"]+`"`)
				authKeyUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "setup_complete") && updates["setup_complete"] != "" {
				updatedLines = append(updatedLines, `setup_complete = `+updates["setup_complete"])
				setupCompleteUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "setup_password") {
				if updates["setup_password"] != "" {
					updatedLines = append(updatedLines, `setup_password = "`+updates["setup_password"]+`"`)
				}
				// If setup_password is empty string, we skip the line (clearing it)
				setupPasswordUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "gh_username") && updates["gh_username"] != "" {
				updatedLines = append(updatedLines, `gh_username = "`+updates["gh_username"]+`"`)
				ghUsernameUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "gh_access_token") && updates["gh_access_token"] != "" {
				updatedLines = append(updatedLines, `gh_access_token = "`+updates["gh_access_token"]+`"`)
				ghAccessTokenUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_host") && updates["smtp_host"] != "" {
				updatedLines = append(updatedLines, `smtp_host = "`+updates["smtp_host"]+`"`)
				smtpHostUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_port") && updates["smtp_port"] != "" {
				updatedLines = append(updatedLines, `smtp_port = "`+updates["smtp_port"]+`"`)
				smtpPortUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_user") && updates["smtp_user"] != "" {
				updatedLines = append(updatedLines, `smtp_user = "`+updates["smtp_user"]+`"`)
				smtpUserUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_pass") && updates["smtp_pass"] != "" {
				updatedLines = append(updatedLines, `smtp_pass = "`+updates["smtp_pass"]+`"`)
				smtpPassUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_from") && !strings.HasPrefix(trimmed, "smtp_from_name") && updates["smtp_from"] != "" {
				updatedLines = append(updatedLines, `smtp_from = "`+updates["smtp_from"]+`"`)
				smtpFromUpdated = true
				continue
			}
			if strings.HasPrefix(trimmed, "smtp_from_name") && updates["smtp_from_name"] != "" {
				updatedLines = append(updatedLines, `smtp_from_name = "`+updates["smtp_from_name"]+`"`)
				smtpFromNameUpdated = true
				continue
			}
		}

		updatedLines = append(updatedLines, line)
	}

	// If still in constants section at the end, add missing fields
	if inConstantsSection {
		if !portUpdated && updates["port"] != "" {
			updatedLines = append(updatedLines, `port = "`+updates["port"]+`"`)
		}
		if !authKeyUpdated && updates["auth_key"] != "" {
			updatedLines = append(updatedLines, `auth_key = "`+updates["auth_key"]+`"`)
		}
		if !setupCompleteUpdated && updates["setup_complete"] != "" {
			updatedLines = append(updatedLines, `setup_complete = `+updates["setup_complete"])
		}
		if !setupPasswordUpdated && updates["setup_password"] != "" {
			updatedLines = append(updatedLines, `setup_password = "`+updates["setup_password"]+`"`)
		}
		if !ghUsernameUpdated && updates["gh_username"] != "" {
			updatedLines = append(updatedLines, `gh_username = "`+updates["gh_username"]+`"`)
		}
		if !ghAccessTokenUpdated && updates["gh_access_token"] != "" {
			updatedLines = append(updatedLines, `gh_access_token = "`+updates["gh_access_token"]+`"`)
		}
		if !smtpHostUpdated && updates["smtp_host"] != "" {
			updatedLines = append(updatedLines, `smtp_host = "`+updates["smtp_host"]+`"`)
		}
		if !smtpPortUpdated && updates["smtp_port"] != "" {
			updatedLines = append(updatedLines, `smtp_port = "`+updates["smtp_port"]+`"`)
		}
		if !smtpUserUpdated && updates["smtp_user"] != "" {
			updatedLines = append(updatedLines, `smtp_user = "`+updates["smtp_user"]+`"`)
		}
		if !smtpPassUpdated && updates["smtp_pass"] != "" {
			updatedLines = append(updatedLines, `smtp_pass = "`+updates["smtp_pass"]+`"`)
		}
		if !smtpFromUpdated && updates["smtp_from"] != "" {
			updatedLines = append(updatedLines, `smtp_from = "`+updates["smtp_from"]+`"`)
		}
		if !smtpFromNameUpdated && updates["smtp_from_name"] != "" {
			updatedLines = append(updatedLines, `smtp_from_name = "`+updates["smtp_from_name"]+`"`)
		}
	} else if !inConstantsSection && len(updates) > 0 {
		// No constants section found, add it at the end
		updatedLines = append(updatedLines, "")
		updatedLines = append(updatedLines, "[constants]")
		if updates["port"] != "" {
			updatedLines = append(updatedLines, `port = "`+updates["port"]+`"`)
		}
		if updates["auth_key"] != "" {
			updatedLines = append(updatedLines, `auth_key = "`+updates["auth_key"]+`"`)
		}
		if updates["setup_complete"] != "" {
			updatedLines = append(updatedLines, `setup_complete = `+updates["setup_complete"])
		}
		if updates["gh_username"] != "" {
			updatedLines = append(updatedLines, `gh_username = "`+updates["gh_username"]+`"`)
		}
		if updates["gh_access_token"] != "" {
			updatedLines = append(updatedLines, `gh_access_token = "`+updates["gh_access_token"]+`"`)
		}
		if updates["smtp_host"] != "" {
			updatedLines = append(updatedLines, `smtp_host = "`+updates["smtp_host"]+`"`)
		}
		if updates["smtp_port"] != "" {
			updatedLines = append(updatedLines, `smtp_port = "`+updates["smtp_port"]+`"`)
		}
		if updates["smtp_user"] != "" {
			updatedLines = append(updatedLines, `smtp_user = "`+updates["smtp_user"]+`"`)
		}
		if updates["smtp_pass"] != "" {
			updatedLines = append(updatedLines, `smtp_pass = "`+updates["smtp_pass"]+`"`)
		}
		if updates["smtp_from"] != "" {
			updatedLines = append(updatedLines, `smtp_from = "`+updates["smtp_from"]+`"`)
		}
		if updates["smtp_from_name"] != "" {
			updatedLines = append(updatedLines, `smtp_from_name = "`+updates["smtp_from_name"]+`"`)
		}
	}

	// Write back to file
	newContent := strings.Join(updatedLines, "\n")
	return os.WriteFile(config_path, []byte(newContent), 0644)
}
