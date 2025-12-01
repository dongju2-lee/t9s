package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Status represents the Git status of a directory
type Status struct {
	Branch         string
	IsDirty        bool
	AheadBy        int
	BehindBy       int
	ModifiedFiles  []string
	UntrackedFiles []string
}

// Manager handles Git operations
type Manager struct{}

// NewManager creates a new Git manager
func NewManager() *Manager {
	return &Manager{}
}

// GetStatus gets the Git status for a directory
func (m *Manager) GetStatus(path string) (*Status, error) {
	status := &Status{}

	// Get current branch
	branch, err := m.getCurrentBranch(path)
	if err != nil {
		return nil, err
	}
	status.Branch = branch

	// Check if working directory is dirty
	isDirty, err := m.isWorkingTreeDirty(path)
	if err != nil {
		return nil, err
	}
	status.IsDirty = isDirty

	// Get modified files
	modifiedFiles, err := m.getModifiedFiles(path)
	if err != nil {
		return nil, err
	}
	status.ModifiedFiles = modifiedFiles

	return status, nil
}

// getCurrentBranch gets the current Git branch
func (m *Manager) getCurrentBranch(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// isWorkingTreeDirty checks if there are uncommitted changes
func (m *Manager) isWorkingTreeDirty(path string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(output) > 0, nil
}

// getModifiedFiles gets the list of modified files
func (m *Manager) getModifiedFiles(path string) ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get modified files: %w", err)
	}

	var files []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Extract filename (skip the first 3 characters which are status codes)
		if len(line) > 3 {
			files = append(files, line[3:])
		}
	}

	return files, nil
}

// GetDiff gets the diff for a specific file or the entire directory
func (m *Manager) GetDiff(path string, file string) (string, error) {
	args := []string{"diff"}
	if file != "" {
		args = append(args, file)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}

	return string(output), nil
}

// GetLastCommit gets information about the last commit
func (m *Manager) GetLastCommit(path string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%h - %s (%an, %ar)")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit: %w", err)
	}

	return string(output), nil
}
