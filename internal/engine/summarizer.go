// Package engine implements commit summary generation for Dikogit.
package engine

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Version is the Dikogit version string.
const Version = "2.0.0"

// FileChange represents a single file changed in a commit.
type FileChange struct {
	Path      string
	Status    string // "added", "modified", "deleted", "renamed", "copied"
	Code      string // "A", "M", "D", "R", "C"
	Additions int
	Deletions int
}

// CommitSummary holds a generated commit title and body.
type CommitSummary struct {
	Title       string
	Body        string
	FullSummary string
}

// Summarizer generates structured conventional commit messages from git diff data.
type Summarizer struct {
	git *Git
}

// NewSummarizer creates a new Summarizer backed by the given Git instance.
func NewSummarizer(g *Git) *Summarizer {
	return &Summarizer{git: g}
}

// GenerateSummary analyzes staged changes and returns a CommitSummary.
// Returns nil if there are no staged changes.
func (s *Summarizer) GenerateSummary() (*CommitSummary, error) {
	changes, err := s.parseChanges()
	if err != nil {
		return nil, err
	}
	if len(changes.Files) == 0 {
		return nil, nil
	}

	commitType := s.determineCommitType(changes)
	totalFiles := len(changes.Files)
	adds := changes.Additions
	dels := changes.Deletions

	// Build scope from top-level directories
	scopeSet := map[string]struct{}{}
	for _, f := range changes.Files {
		parts := strings.Split(f.Path, "/")
		if len(parts) > 1 {
			scopeSet[parts[0]] = struct{}{}
		} else {
			scopeSet["root"] = struct{}{}
		}
	}
	scopes := sortedKeys(scopeSet)
	if len(scopes) > 2 {
		scopes = scopes[:2]
	}
	scopeStr := ""
	if len(scopes) > 0 {
		scopeStr = "(" + strings.Join(scopes, ", ") + ")"
	}

	// Build title
	var title string
	if totalFiles == 1 {
		base := filepath.Base(changes.Files[0].Path)
		title = fmt.Sprintf("%s%s: update %s", commitType, scopeStr, base)
	} else {
		title = fmt.Sprintf("%s%s: update %d files (+%d -%d)", commitType, scopeStr, totalFiles, adds, dels)
	}

	// Build body
	var bodyLines []string
	bodyLines = append(bodyLines, "### Dikogits Commit Summary")
	bodyLines = append(bodyLines, fmt.Sprintf("- **Stats**: %d file(s) changed, %d insertion(s)(+), %d deletion(s)(-)", totalFiles, adds, dels))
	bodyLines = append(bodyLines, "- **File Breakdown**:")

	limit := 10
	if len(changes.Files) < limit {
		limit = len(changes.Files)
	}
	for _, f := range changes.Files[:limit] {
		sym := statusSymbol(f.Status)
		bodyLines = append(bodyLines, fmt.Sprintf("  - %s `%s` (+%d -%d)", sym, f.Path, f.Additions, f.Deletions))
	}
	if len(changes.Files) > 10 {
		bodyLines = append(bodyLines, fmt.Sprintf("  - ... and %d more files", len(changes.Files)-10))
	}

	if len(changes.Details) > 0 {
		bodyLines = append(bodyLines, "- **Key Changes**:")
		for _, d := range changes.Details {
			bodyLines = append(bodyLines, "  - "+d)
		}
	}

	body := strings.Join(bodyLines, "\n")
	full := title + "\n\n" + body

	return &CommitSummary{
		Title:       title,
		Body:        body,
		FullSummary: full,
	}, nil
}

// changeSet is the internal structured result of parsing staged diffs.
type changeSet struct {
	Files     []FileChange
	Additions int
	Deletions int
	Details   []string
}

func (s *Summarizer) parseChanges() (*changeSet, error) {
	nameStatus := s.git.StagedNameStatus()
	numstat := s.git.StagedNumstat()
	diffText := s.git.StagedDiff()

	// Parse numstat -> map path -> (adds, dels)
	numstatMap := map[string][2]int{}
	addTotal, delTotal := 0, 0
	for _, line := range splitLines(numstat) {
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		adds := parseIntSafe(parts[0])
		dels := parseIntSafe(parts[1])
		path := parts[2]
		numstatMap[path] = [2]int{adds, dels}
		addTotal += adds
		delTotal += dels
	}

	// Parse name-status -> file list
	var files []FileChange
	for _, line := range splitLines(nameStatus) {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}
		code := string(parts[0][0])
		path := parts[len(parts)-1]

		statusName := statusName(code)
		nums := numstatMap[path]

		files = append(files, FileChange{
			Path:      path,
			Status:    statusName,
			Code:      code,
			Additions: nums[0],
			Deletions: nums[1],
		})
	}

	details := extractDiffDetails(diffText)

	return &changeSet{
		Files:     files,
		Additions: addTotal,
		Deletions: delTotal,
		Details:   details,
	}, nil
}

// determineCommitType classifies the change as a conventional commit type.
func (s *Summarizer) determineCommitType(cs *changeSet) string {
	var paths, statuses []string
	for _, f := range cs.Files {
		paths = append(paths, strings.ToLower(f.Path))
		statuses = append(statuses, f.Status)
	}

	// Test files
	if allMatch(paths, func(p string) bool {
		return strings.Contains(p, "test") || strings.Contains(p, "spec") || strings.HasPrefix(p, "tests/")
	}) {
		return "test"
	}

	// Documentation
	if allMatch(paths, func(p string) bool {
		return hasSuffix(p, ".md", ".rst", ".txt") || strings.Contains(p, "docs/") || strings.Contains(p, "readme")
	}) {
		return "docs"
	}

	// Config / CI
	if allMatch(paths, func(p string) bool {
		return hasSuffix(p, ".yml", ".yaml", ".json", ".toml", ".ini") ||
			strings.Contains(p, ".github/") || strings.Contains(p, "docker")
	}) {
		return "chore"
	}

	// Fix indicators in diff
	diff := strings.ToLower(s.git.StagedDiff())
	if strings.Contains(diff, "fix") || strings.Contains(diff, "bug") ||
		strings.Contains(diff, "error") || strings.Contains(diff, "patch") {
		return "fix"
	}

	// New files → feature
	for _, st := range statuses {
		if st == "added" {
			return "feat"
		}
	}

	// Modified files → refactor
	for _, st := range statuses {
		if st == "modified" {
			return "refactor"
		}
	}

	return "chore"
}

// extractDiffDetails extracts notable lines from a unified diff (functions, imports, headings).
func extractDiffDetails(diff string) []string {
	keywordRe := regexp.MustCompile(`^(def |class |func |function |async def |export |pub fn |interface |type |const |let |var )`)
	importRe := regexp.MustCompile(`^(import |from |require\(|#include|use )`)

	seen := map[string]struct{}{}
	var results []string

	for _, line := range strings.Split(diff, "\n") {
		if !strings.HasPrefix(line, "+") || strings.HasPrefix(line, "+++") {
			continue
		}
		content := strings.TrimSpace(line[1:])

		var detail string
		switch {
		case keywordRe.MatchString(content):
			if len(content) > 60 {
				content = content[:60]
			}
			detail = fmt.Sprintf("Defined `%s`", content)
		case importRe.MatchString(content):
			if len(content) > 60 {
				content = content[:60]
			}
			detail = fmt.Sprintf("Added import: `%s`", content)
		case strings.HasPrefix(content, "#") && len(content) < 80:
			detail = fmt.Sprintf("Added section: `%s`", content)
		}

		if detail != "" {
			if _, ok := seen[detail]; !ok {
				seen[detail] = struct{}{}
				results = append(results, detail)
				if len(results) == 5 {
					break
				}
			}
		}
	}
	return results
}

// --- helpers ---

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func parseIntSafe(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return n
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func statusName(code string) string {
	switch code {
	case "A":
		return "added"
	case "M":
		return "modified"
	case "D":
		return "deleted"
	case "R":
		return "renamed"
	case "C":
		return "copied"
	default:
		return "changed"
	}
}

func statusSymbol(status string) string {
	switch status {
	case "added":
		return "[A]"
	case "modified":
		return "[M]"
	case "deleted":
		return "[D]"
	case "renamed":
		return "[R]"
	case "copied":
		return "[C]"
	default:
		return "[*]"
	}
}

func allMatch(items []string, fn func(string) bool) bool {
	if len(items) == 0 {
		return false
	}
	for _, it := range items {
		if !fn(it) {
			return false
		}
	}
	return true
}

func hasSuffix(s string, suffixes ...string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
