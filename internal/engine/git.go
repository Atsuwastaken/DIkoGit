// Package engine provides the core git-shelling logic for Dikogit.
// It shells out to the system git binary for all git operations.
package engine

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Git wraps shell calls to the system `git` binary.
type Git struct {
	RepoRoot string
}

// NewGit creates a Git instance rooted at the current working directory's repo root.
// Returns an error if the working directory is not inside a git repository.
func NewGit() (*Git, error) {
	root, err := runGit("", "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("fatal: not a git repository (or any of the parent directories): .git")
	}
	return &Git{RepoRoot: strings.TrimSpace(root)}, nil
}

// IsClean returns true if there are no changes (staged or unstaged) in the working tree.
func (g *Git) IsClean() bool {
	out, err := runGit(g.RepoRoot, "status", "--porcelain")
	if err != nil {
		return true
	}
	return strings.TrimSpace(out) == ""
}

// HasStagedChanges returns true if there are staged changes.
func (g *Git) HasStagedChanges() bool {
	out, err := runGit(g.RepoRoot, "diff", "--cached", "--name-only")
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) != ""
}

// AddAll stages all changes (equivalent to `git add .`).
func (g *Git) AddAll() error {
	_, err := runGit(g.RepoRoot, "add", ".")
	if err != nil {
		return fmt.Errorf("git add . failed: %w", err)
	}
	return nil
}

// Commit creates a commit with the given title and body.
func (g *Git) Commit(title, body string) (string, error) {
	out, err := runGit(g.RepoRoot, "commit", "-m", title, "-m", body)
	if err != nil {
		return "", fmt.Errorf("git commit failed: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// Push pushes to the tracked remote of the current branch.
// Returns a descriptive error if no remote is configured.
func (g *Git) Push() (string, error) {
	// Check if any remote exists
	remotes, _ := runGit(g.RepoRoot, "remote")
	if strings.TrimSpace(remotes) == "" {
		return "", fmt.Errorf("no remote configured: your commit is saved locally")
	}

	out, err := runGit(g.RepoRoot, "push")
	if err != nil {
		// Try push with upstream set if no tracking branch
		if strings.Contains(err.Error(), "no upstream") || strings.Contains(err.Error(), "set-upstream") {
			branch, _ := g.CurrentBranch()
			out, err = runGit(g.RepoRoot, "push", "--set-upstream", "origin", branch)
			if err != nil {
				return "", fmt.Errorf("git push failed: %w", err)
			}
			return strings.TrimSpace(out), nil
		}
		return "", fmt.Errorf("git push failed: %w", err)
	}
	return strings.TrimSpace(out), nil
}

// CurrentBranch returns the name of the currently checked-out branch.
func (g *Git) CurrentBranch() (string, error) {
	out, err := runGit(g.RepoRoot, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// StagedNameStatus returns the `--name-status` output for staged changes.
func (g *Git) StagedNameStatus() string {
	out, _ := runGit(g.RepoRoot, "diff", "--cached", "--name-status")
	return strings.TrimSpace(out)
}

// StagedNumstat returns the `--numstat` output for staged changes.
func (g *Git) StagedNumstat() string {
	out, _ := runGit(g.RepoRoot, "diff", "--cached", "--numstat")
	return strings.TrimSpace(out)
}

// StagedDiff returns the full unified diff for staged changes.
func (g *Git) StagedDiff() string {
	out, _ := runGit(g.RepoRoot, "diff", "--cached")
	return out
}

// runGit runs git with the given args in the given directory (or cwd if dir is "").
// It returns stdout on success, or a combined error with stderr on failure.
func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = strings.TrimSpace(stdout.String())
		}
		return "", fmt.Errorf("%s", errMsg)
	}
	return stdout.String(), nil
}
