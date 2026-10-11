// Dikogit - Fully automated Git CLI extension.
//
// Usage:
//
//	git dikogit              — stage all, commit with generated summary, then push
//	git dikogit generate     — print commit summary for staged changes (no commit)
//	git dikogit generate --json  — same but in JSON format
//	git dikogit --version    — print version
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Atsuwastaken/dikogit/internal/engine"
)

func main() {
	args := os.Args[1:]

	// --version / -v
	if hasFlag(args, "--version", "-v") {
		fmt.Printf("git-dikogit %s\n", engine.Version)
		os.Exit(0)
	}

	// Resolve subcommand
	subcommand := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		subcommand = args[0]
		args = args[1:]
	}

	// Initialise git engine — validates we are inside a repo
	g, err := engine.NewGit()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	summarizer := engine.NewSummarizer(g)

	switch subcommand {
	case "", "commit":
		runAutoCommit(g, summarizer)
	case "generate":
		jsonOut := hasFlag(args, "--json")
		runGenerate(g, summarizer, jsonOut)
	default:
		printHelp()
		os.Exit(1)
	}
}

// runAutoCommit performs: git add . → commit with summary → git push
func runAutoCommit(g *engine.Git, s *engine.Summarizer) {
	// 1. Check for any changes at all
	if g.IsClean() {
		fmt.Fprintln(os.Stderr, "No changes found to commit (working tree clean).")
		os.Exit(1)
	}

	// 2. git add .
	if err := g.AddAll(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to stage changes: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Staged all changes (git add .)")

	// 3. Generate summary from staged changes
	summary, err := s.GenerateSummary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to generate summary: %v\n", err)
		os.Exit(1)
	}
	if summary == nil {
		fmt.Fprintln(os.Stderr, "No staged changes found to commit.")
		os.Exit(1)
	}

	// 4. git commit
	out, err := g.Commit(summary.Title, summary.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Commit failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Committed with summary:")
	fmt.Println()
	fmt.Println(summary.FullSummary)
	if out != "" {
		fmt.Println()
		fmt.Println(out)
	}

	// 5. git push
	fmt.Println()
	pushOut, err := g.Push()
	if err != nil {
		// Warn but do NOT fail — commit is already saved locally
		fmt.Fprintf(os.Stderr, "⚠  Push skipped: %v\n", err)
		fmt.Fprintln(os.Stderr, "   Your commit is saved locally. Run `git push` manually when ready.")
		os.Exit(0)
	}
	fmt.Println("✓ Pushed to remote")
	if pushOut != "" {
		fmt.Println(pushOut)
	}
}

// runGenerate prints the commit summary for staged changes without committing.
func runGenerate(g *engine.Git, s *engine.Summarizer, asJSON bool) {
	if !g.HasStagedChanges() {
		fmt.Fprintln(os.Stderr, "No staged changes found. Stage files first with `git add`.")
		os.Exit(1)
	}

	summary, err := s.GenerateSummary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ Failed to generate summary: %v\n", err)
		os.Exit(1)
	}
	if summary == nil {
		fmt.Fprintln(os.Stderr, "No staged changes found.")
		os.Exit(1)
	}

	if asJSON {
		out, _ := json.MarshalIndent(map[string]string{
			"title":        summary.Title,
			"body":         summary.Body,
			"full_summary": summary.FullSummary,
		}, "", "  ")
		fmt.Println(string(out))
	} else {
		fmt.Println(summary.FullSummary)
	}
}

func printHelp() {
	fmt.Printf(`git-dikogit %s — Fully automated Git commit & push

USAGE:
  git dikogit                    Stage all, commit with auto-generated summary, and push
  git dikogit generate           Print commit summary for staged changes (no commit)
  git dikogit generate --json    Same as above, in JSON format
  git dikogit --version          Print version

EXAMPLES:
  git dikogit                    # full automation: add → commit → push
  git add .
  git dikogit generate           # preview the commit message
  git commit -m "$(git dikogit generate)"
`, engine.Version)
}

func hasFlag(args []string, flags ...string) bool {
	for _, a := range args {
		for _, f := range flags {
			if a == f {
				return true
			}
		}
	}
	return false
}
