package dao

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/idongju/t9s/internal/model"
)

// GitDAO handles Git data access operations
type GitDAO struct{}

// NewGitDAO creates a new Git DAO
func NewGitDAO() *GitDAO {
	return &GitDAO{}
}

// GetStatus gets the Git status for a directory
func (d *GitDAO) GetStatus(path string) (*model.GitStatus, error) {
	status := &model.GitStatus{}

	// Get current branch
	branch, err := d.getCurrentBranch(path)
	if err != nil {
		return nil, err
	}
	status.Branch = branch

	// Check if working directory is dirty
	isDirty, err := d.isWorkingTreeDirty(path)
	if err != nil {
		return nil, err
	}
	status.IsDirty = isDirty

	// Get modified files
	modifiedFiles, err := d.getModifiedFiles(path)
	if err != nil {
		return nil, err
	}
	status.ModifiedFiles = modifiedFiles

	// Get last commit info
	lastCommit, err := d.getLastCommit(path)
	if err == nil {
		status.LastCommit = lastCommit
	}

	return status, nil
}

// getCurrentBranch gets the current Git branch
func (d *GitDAO) getCurrentBranch(path string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// isWorkingTreeDirty checks if there are uncommitted changes
func (d *GitDAO) isWorkingTreeDirty(path string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(output) > 0, nil
}

// getModifiedFiles gets the list of modified files
func (d *GitDAO) getModifiedFiles(path string) ([]string, error) {
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
		if len(line) > 3 {
			files = append(files, line[3:])
		}
	}

	return files, nil
}

// GetDiff gets the diff for a specific file or the entire directory
func (d *GitDAO) GetDiff(path string, file string) (string, error) {
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

// getLastCommit gets information about the last commit
func (d *GitDAO) getLastCommit(path string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%h - %s (%an, %ar)")
	cmd.Dir = path

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit: %w", err)
	}

	return string(output), nil
}

