package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Directory represents a Terraform directory with its metadata
type Directory struct {
	Name         string
	Path         string
	Status       Status
	LastApply    time.Time
	ConfigPath   string
	TfvarsFiles  []string
	BackendType  string
	BackendKey   string
	HasDrift     bool
	GitStatus    string
}

// Status represents the current state of a Terraform directory
type Status int

const (
	StatusUnknown Status = iota
	StatusSynced
	StatusDrift
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusSynced:
		return "✓ Synced"
	case StatusDrift:
		return "⚠ Drift"
	case StatusError:
		return "✗ Error"
	default:
		return "? Unknown"
	}
}

// Manager handles Terraform operations
type Manager struct {
	RootPath string
}

// NewManager creates a new Terraform manager
func NewManager(rootPath string) *Manager {
	return &Manager{
		RootPath: rootPath,
	}
}

// ScanDirectories scans for Terraform directories
func (m *Manager) ScanDirectories() ([]*Directory, error) {
	var directories []*Directory

	entries, err := os.ReadDir(m.RootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read root directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(m.RootPath, entry.Name())
		
		// Check if it's a Terraform directory (has .tf files)
		if !m.isTerraformDir(dirPath) {
			continue
		}

		dir := &Directory{
			Name:   entry.Name(),
			Path:   dirPath,
			Status: StatusUnknown,
		}

		// Find config directory and tfvars files
		configPath := filepath.Join(dirPath, "config")
		if stat, err := os.Stat(configPath); err == nil && stat.IsDir() {
			dir.ConfigPath = configPath
			dir.TfvarsFiles = m.findTfvarsFiles(configPath)
		}

		// Get backend configuration
		m.parseBackendConfig(dir)

		directories = append(directories, dir)
	}

	return directories, nil
}

// isTerraformDir checks if a directory contains Terraform files
func (m *Manager) isTerraformDir(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tf") {
			return true
		}
	}
	return false
}

// findTfvarsFiles finds all .tfvars files in a directory
func (m *Manager) findTfvarsFiles(path string) []string {
	var tfvarsFiles []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return tfvarsFiles
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tfvars") {
			tfvarsFiles = append(tfvarsFiles, entry.Name())
		}
	}

	return tfvarsFiles
}

// parseBackendConfig parses the backend configuration from Terraform files
func (m *Manager) parseBackendConfig(dir *Directory) {
	// This is a simplified version - in production, you'd parse the actual .tf files
	// For now, we'll check for a terraform.tfstate file or .terraform directory
	
	terraformDir := filepath.Join(dir.Path, ".terraform")
	if stat, err := os.Stat(terraformDir); err == nil && stat.IsDir() {
		dir.BackendType = "s3" // Assuming S3, should parse from backend.tf
		dir.BackendKey = fmt.Sprintf("terraform-state/%s/terraform.tfstate", dir.Name)
	}
}

// CheckDrift checks if there is drift between state and actual infrastructure
func (m *Manager) CheckDrift(dir *Directory) error {
	cmd := exec.Command("terraform", "plan", "-detailed-exitcode")
	cmd.Dir = dir.Path

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Exit code 2 means there are changes (drift detected)
			if exitErr.ExitCode() == 2 {
				dir.HasDrift = true
				dir.Status = StatusDrift
				return nil
			}
		}
		return fmt.Errorf("terraform plan failed: %w: %s", err, string(output))
	}

	// Exit code 0 means no changes
	dir.HasDrift = false
	dir.Status = StatusSynced
	return nil
}

// GetStateInfo retrieves state information
func (m *Manager) GetStateInfo(dir *Directory) (map[string]interface{}, error) {
	cmd := exec.Command("terraform", "show", "-json")
	cmd.Dir = dir.Path

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("terraform show failed: %w", err)
	}

	var state map[string]interface{}
	if err := json.Unmarshal(output, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}

	return state, nil
}

// Plan runs terraform plan
func (m *Manager) Plan(dir *Directory, tfvarsFile string) (string, error) {
	args := []string{"plan"}
	if tfvarsFile != "" {
		args = append(args, "-var-file="+filepath.Join(dir.ConfigPath, tfvarsFile))
	}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir.Path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("terraform plan failed: %w", err)
	}

	return string(output), nil
}

// Apply runs terraform apply
func (m *Manager) Apply(dir *Directory, tfvarsFile string) (string, error) {
	args := []string{"apply", "-auto-approve"}
	if tfvarsFile != "" {
		args = append(args, "-var-file="+filepath.Join(dir.ConfigPath, tfvarsFile))
	}

	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir.Path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("terraform apply failed: %w", err)
	}

	dir.LastApply = time.Now()
	return string(output), nil
}

// GetHelmReleases extracts Helm releases from Terraform state
func (m *Manager) GetHelmReleases(dir *Directory) ([]HelmRelease, error) {
	stateInfo, err := m.GetStateInfo(dir)
	if err != nil {
		return nil, err
	}

	var releases []HelmRelease
	
	// Parse state for helm_release resources
	if values, ok := stateInfo["values"].(map[string]interface{}); ok {
		if rootModule, ok := values["root_module"].(map[string]interface{}); ok {
			if resources, ok := rootModule["resources"].([]interface{}); ok {
				for _, res := range resources {
					if resource, ok := res.(map[string]interface{}); ok {
						if resType, ok := resource["type"].(string); ok && resType == "helm_release" {
							release := HelmRelease{
								Name: resource["name"].(string),
							}
							if values, ok := resource["values"].(map[string]interface{}); ok {
								if chart, ok := values["chart"].(string); ok {
									release.Chart = chart
								}
								if version, ok := values["version"].(string); ok {
									release.Version = version
								}
							}
							releases = append(releases, release)
						}
					}
				}
			}
		}
	}

	return releases, nil
}

// HelmRelease represents a Helm release managed by Terraform
type HelmRelease struct {
	Name      string
	Chart     string
	Version   string
	Namespace string
}
