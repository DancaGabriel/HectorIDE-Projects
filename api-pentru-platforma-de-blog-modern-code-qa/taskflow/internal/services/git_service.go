package main

import (
	"fmt"
	"os"
	"os/exec"
)

// GitService provides git-related functionalities
type GitService struct {
	repoPath string
}

// NewGitService creates a new GitService instance.
func NewGitService(repoPath string) *GitService {
	return &GitService{repoPath: repoPath}
}

// InitializeGitRepo initializes a new Git repository if one doesn't exist.
func (g *GitService) InitializeGitRepo() error {
	if _, err := os.Stat(g.repoPath + "/.git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", g.repoPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %w", err)
		}
		fmt.Println("Initialized new Git repository at", g.repoPath)
	} else if err != nil {
		return fmt.Errorf("error checking git repository: %w", err)
	}
	return nil
}

// CommitAndPush commits changes and pushes them to a remote repository (if configured).
func (g *GitService) CommitAndPush(message string) error {
	cmdAdd := exec.Command("git", "-C", g.repoPath, "add", ".")
	if err := cmdAdd.Run(); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	cmdCommit := exec.Command("git", "-C", g.repoPath, "commit", "-m", message)
	if err := cmdCommit.Run(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	// Attempt to push changes.  This might fail if no remote is configured.
	cmdPush := exec.Command("git", "-C", g.repoPath, "push")
	if err := cmdPush.Run(); err != nil {
		fmt.Println("Warning: Failed to push changes (no remote configured or other issues):", err)
		// We don't want to fail the entire operation if the push fails.
		// Pushing to remote could be an optional feature.
	}
	return nil
}