package auxiliary

import (
	"runtime"

	"github.com/laurent22/toml-go"
)

/*
Get required field from config file.
*/
func GetFromConfig(configField string) string {
	os := runtime.GOOS
	config_path := "/csmkactionhelper/config/config.toml"
	var parser toml.Parser
	if os == "windows" {
		config_path = "C:\\csmkactionhelper\\config\\config.toml"
	}

	config := parser.ParseFile(config_path)
	return config.GetString(configField)
}
