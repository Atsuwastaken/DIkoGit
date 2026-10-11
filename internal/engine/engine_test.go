package engine_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Atsuwastaken/dikogit/internal/engine"
)

// initTestRepo creates a temporary git repository for tests.
func initTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v failed: %s", args, out)
		}
	}
	run("init")
	run("config", "user.name", "Test User")
	run("config", "user.email", "test@example.com")
	return dir
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func gitAdd(t *testing.T, dir string, files ...string) {
	t.Helper()
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git add failed: %s", out)
	}
}

// TestNotAGitRepo verifies NewGit fails outside a repo.
func TestNotAGitRepo(t *testing.T) {
	dir := t.TempDir() // not a git repo
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	_, err := engine.NewGit()
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

// TestIsClean verifies a freshly initialised repo is clean.
func TestIsClean(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	g, err := engine.NewGit()
	if err != nil {
		t.Fatal(err)
	}
	if !g.IsClean() {
		t.Fatal("expected clean repo")
	}
}

// TestIsCleanWithChanges verifies IsClean returns false after writing an untracked file.
func TestIsCleanWithChanges(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	writeFile(t, dir, "hello.go", "package main\n")

	g, err := engine.NewGit()
	if err != nil {
		t.Fatal(err)
	}
	if g.IsClean() {
		t.Fatal("expected dirty repo after adding file")
	}
}

// TestAddAll verifies git add . stages all files.
func TestAddAll(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	writeFile(t, dir, "app.go", "package main\nfunc main(){}\n")

	g, err := engine.NewGit()
	if err != nil {
		t.Fatal(err)
	}
	if err := g.AddAll(); err != nil {
		t.Fatalf("AddAll failed: %v", err)
	}
	if !g.HasStagedChanges() {
		t.Fatal("expected staged changes after AddAll")
	}
}

// TestGenerateSummaryFeat verifies a feature summary is classified correctly.
func TestGenerateSummaryFeat(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	writeFile(t, dir, "feature.go", "package main\nfunc newFeature() {}\n")
	gitAdd(t, dir, "feature.go")

	g, _ := engine.NewGit()
	s := engine.NewSummarizer(g)

	summary, err := s.GenerateSummary()
	if err != nil {
		t.Fatal(err)
	}
	if summary == nil {
		t.Fatal("expected non-nil summary")
	}
	if !strings.HasPrefix(summary.Title, "feat") {
		t.Fatalf("expected title to start with 'feat', got: %s", summary.Title)
	}
	if !strings.Contains(summary.Body, "feature.go") {
		t.Fatalf("expected body to contain 'feature.go', got: %s", summary.Body)
	}
}

// TestGenerateSummaryDocs verifies a docs summary is classified correctly.
func TestGenerateSummaryDocs(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	writeFile(t, dir, "README.md", "# Docs\nThis is documentation.\n")
	gitAdd(t, dir, "README.md")

	g, _ := engine.NewGit()
	s := engine.NewSummarizer(g)

	summary, err := s.GenerateSummary()
	if err != nil {
		t.Fatal(err)
	}
	if summary == nil {
		t.Fatal("expected non-nil summary")
	}
	if !strings.HasPrefix(summary.Title, "docs") {
		t.Fatalf("expected title to start with 'docs', got: %s", summary.Title)
	}
}

// TestGenerateSummaryNoStaged verifies nil is returned when nothing is staged.
func TestGenerateSummaryNoStaged(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	g, _ := engine.NewGit()
	s := engine.NewSummarizer(g)

	summary, err := s.GenerateSummary()
	if err != nil {
		t.Fatal(err)
	}
	if summary != nil {
		t.Fatalf("expected nil summary for empty staging area, got: %+v", summary)
	}
}

// TestCommit verifies a commit is created with the correct message.
func TestCommit(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	writeFile(t, dir, "main.go", "package main\nfunc main(){}\n")
	gitAdd(t, dir, "main.go")

	g, _ := engine.NewGit()
	s := engine.NewSummarizer(g)

	summary, err := s.GenerateSummary()
	if err != nil || summary == nil {
		t.Fatal("unexpected summary failure")
	}

	out, err := g.Commit(summary.Title, summary.Body)
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
	_ = out

	// Verify commit log
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = dir
	logOut, _ := cmd.Output()
	if !strings.Contains(string(logOut), summary.Title) {
		t.Fatalf("commit log missing title, got: %s", logOut)
	}
	if !strings.Contains(string(logOut), "Dikogits Commit Summary") {
		t.Fatalf("commit log missing body, got: %s", logOut)
	}
}

// TestFullAutoCommitFlow verifies IsClean → AddAll → GenerateSummary → Commit flow end-to-end.
func TestFullAutoCommitFlow(t *testing.T) {
	dir := initTestRepo(t)
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	writeFile(t, dir, "calculator.go", "package main\nfunc add(a,b int) int { return a+b }\n")

	g, err := engine.NewGit()
	if err != nil {
		t.Fatal(err)
	}
	if g.IsClean() {
		t.Fatal("expected dirty repo")
	}

	if err := g.AddAll(); err != nil {
		t.Fatal(err)
	}

	s := engine.NewSummarizer(g)
	summary, err := s.GenerateSummary()
	if err != nil || summary == nil {
		t.Fatal("summary failed")
	}

	out, err := g.Commit(summary.Title, summary.Body)
	if err != nil {
		t.Fatalf("commit failed: %v", err)
	}
	_ = out

	if !g.IsClean() {
		t.Fatal("expected clean repo after commit")
	}
}

// TestVersion checks the version constant is set.
func TestVersion(t *testing.T) {
	if engine.Version == "" {
		t.Fatal("engine.Version must not be empty")
	}
}
