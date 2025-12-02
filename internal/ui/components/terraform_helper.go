package components

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/idongju/t9s/internal/config"
)

// TerraformCommandInfo holds information about a terraform command
type TerraformCommandInfo struct {
	WorkDir    string
	ConfigFile string
	Content    string
	Command    string
}

// GetTerraformCommandInfo prepares terraform command information
func GetTerraformCommandInfo(path string, template string, configFileName string, cfg *config.Config) *TerraformCommandInfo {
	info := &TerraformCommandInfo{}

	// Determine working directory
	fileInfo, err := os.Stat(path)
	if err == nil && !fileInfo.IsDir() {
		info.WorkDir = filepath.Dir(path)
		// If it's a config file, use it directly
		if strings.Contains(path, configFileName) {
			info.ConfigFile = path
		}
	} else {
		info.WorkDir = path
	}

	// Find config file if not already set
	if info.ConfigFile == "" {
		configPath := filepath.Join(info.WorkDir, configFileName)
		if _, err := os.Stat(configPath); err == nil {
			info.ConfigFile = configPath
		}
	}

	// Read config file content
	if info.ConfigFile != "" {
		content, err := ioutil.ReadFile(info.ConfigFile)
		if err == nil {
			info.Content = string(content)
		}
	}

	// Build command
	cmdStr := template
	if info.ConfigFile != "" {
		cmdStr = strings.ReplaceAll(cmdStr, "{varfile}", info.ConfigFile)
		cmdStr = strings.ReplaceAll(cmdStr, "{initconf}", info.ConfigFile)
	} else {
		cmdStr = strings.ReplaceAll(cmdStr, "-var-file={varfile}", "")
		cmdStr = strings.ReplaceAll(cmdStr, "-backend-config={initconf}", "")
	}
	info.Command = cmdStr

	return info
}

