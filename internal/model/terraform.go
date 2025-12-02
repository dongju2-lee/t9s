package model

import "time"

// TerraformDirectory represents a Terraform directory with its metadata
type TerraformDirectory struct {
	Name         string
	Path         string
	Status       TerraformStatus
	LastApply    time.Time
	ConfigPath   string
	TfvarsFiles  []string
	BackendType  string
	BackendKey   string
	HasDrift     bool
	GitBranch    string
	GitDirty     bool
	Resources    int
	Outputs      int
}

// TerraformStatus represents the current state of a Terraform directory
type TerraformStatus int

const (
	StatusUnknown TerraformStatus = iota
	StatusSynced
	StatusDrift
	StatusError
	StatusPending
)

func (s TerraformStatus) String() string {
	switch s {
	case StatusSynced:
		return "✓ Synced"
	case StatusDrift:
		return "⚠ Drift"
	case StatusError:
		return "✗ Error"
	case StatusPending:
		return "⌛ Pending"
	default:
		return "? Unknown"
	}
}

// HelmRelease represents a Helm release managed by Terraform
type HelmRelease struct {
	Name      string
	Chart     string
	Version   string
	Namespace string
	Status    string
	Updated   time.Time
}

