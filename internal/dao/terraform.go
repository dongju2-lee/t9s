package dao

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/idongju/t9s/internal/model"
)

// TerraformDAO handles Terraform data access operations
type TerraformDAO struct {
	RootPath string
}

// NewTerraformDAO creates a new Terraform DAO
func NewTerraformDAO(rootPath string) *TerraformDAO {
	return &TerraformDAO{
		RootPath: rootPath,
	}
}

// ListDirectories scans for Terraform directories
func (d *TerraformDAO) ListDirectories() ([]*model.TerraformDirectory, error) {
	var directories []*model.TerraformDirectory

	entries, err := os.ReadDir(d.RootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read root directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(d.RootPath, entry.Name())
		
		// Check if it's a Terraform directory (has .tf files)
		if !d.isTerraformDir(dirPath) {
			continue
		}

		dir := &model.TerraformDirectory{
			Name:   entry.Name(),
			Path:   dirPath,
			Status: model.StatusUnknown,
		}

		// Find config directory and tfvars files
		configPath := filepath.Join(dirPath, "config")
		if stat, err := os.Stat(configPath); err == nil && stat.IsDir() {
			dir.ConfigPath = configPath
			dir.TfvarsFiles = d.findTfvarsFiles(configPath)
		}

		// Get backend configuration
		d.parseBackendConfig(dir)

		directories = append(directories, dir)
	}

	return directories, nil
}

// isTerraformDir checks if a directory contains Terraform files
func (d *TerraformDAO) isTerraformDir(path string) bool {
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
func (d *TerraformDAO) findTfvarsFiles(path string) []string {
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
func (d *TerraformDAO) parseBackendConfig(dir *model.TerraformDirectory) {
	terraformDir := filepath.Join(dir.Path, ".terraform")
	if stat, err := os.Stat(terraformDir); err == nil && stat.IsDir() {
		dir.BackendType = "s3"
		dir.BackendKey = fmt.Sprintf("terraform-state/%s/terraform.tfstate", dir.Name)
	}
}

// CheckDrift checks if there is drift between state and actual infrastructure
func (d *TerraformDAO) CheckDrift(dir *model.TerraformDirectory) error {
	cmd := exec.Command("terraform", "plan", "-detailed-exitcode")
	cmd.Dir = dir.Path

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 2 {
				dir.HasDrift = true
				dir.Status = model.StatusDrift
				return nil
			}
		}
		return fmt.Errorf("terraform plan failed: %w: %s", err, string(output))
	}

	dir.HasDrift = false
	dir.Status = model.StatusSynced
	return nil
}

// GetStateInfo retrieves state information
func (d *TerraformDAO) GetStateInfo(dir *model.TerraformDirectory) (map[string]interface{}, error) {
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
func (d *TerraformDAO) Plan(dir *model.TerraformDirectory, tfvarsFile string) (string, error) {
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
func (d *TerraformDAO) Apply(dir *model.TerraformDirectory, tfvarsFile string) (string, error) {
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
func (d *TerraformDAO) GetHelmReleases(dir *model.TerraformDirectory) ([]*model.HelmRelease, error) {
	stateInfo, err := d.GetStateInfo(dir)
	if err != nil {
		return nil, err
	}

	var releases []*model.HelmRelease
	
	if values, ok := stateInfo["values"].(map[string]interface{}); ok {
		if rootModule, ok := values["root_module"].(map[string]interface{}); ok {
			if resources, ok := rootModule["resources"].([]interface{}); ok {
				for _, res := range resources {
					if resource, ok := res.(map[string]interface{}); ok {
						if resType, ok := resource["type"].(string); ok && resType == "helm_release" {
							release := &model.HelmRelease{
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


