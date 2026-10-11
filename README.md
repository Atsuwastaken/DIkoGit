# Dikogit - Git CLI Extension for macOS

**Dikogit** is a fully automated Git CLI extension written in Go. Run `git dikogit` to automatically stage all changes, generate a structured conventional commit message, commit, and push — all with a single command.

---

## Features

- 🚀 **Full Automation**: `git dikogit` runs `git add .` → `git commit` → `git push` in one shot
- 🏷️ **Conventional Commit Messages**: Auto-classifies commits as `feat`, `fix`, `docs`, `refactor`, `test`, or `chore`
- 📊 **Detailed Commit Body**: Shows file stats, additions/deletions, and key code constructs
- 🍺 **Homebrew Distribution**: Easy global installation via Homebrew tap
- 💻 **Flexible CLI**: Auto-commit or preview summaries before committing
- ⚡ **Zero Runtime Dependencies**: Single static Go binary — no Python, no Node, no runtimes
- 🔒 **Safe Push**: If no remote is configured, warns gracefully and keeps the commit locally

---

## Installation

Install `dikogit` globally on macOS via Homebrew:

```bash
# 1. Tap the custom Homebrew repository
brew tap Atsuwastaken/dikogit

# 2. Install dikogit
brew install dikogit
```

Verify installation:
```bash
git dikogit --version
```

---

## Usage

### 1. Fully Automated: Stage → Commit → Push (Default)

Run from any git repository:

```bash
git dikogit
```

This will:
1. `git add .` — stage all changes
2. Analyze the diff and generate a structured commit message
3. `git commit` — commit with the generated message
4. `git push` — push to the remote

Example output:
```text
✓ Staged all changes (git add .)
✓ Committed with summary:

feat(cmd): update main.go

### Dikogits Commit Summary
- **Stats**: 1 file(s) changed, 42 insertion(s)(+), 3 deletion(s)(-)
- **File Breakdown**:
  - [M] `cmd/git-dikogit/main.go` (+42 -3)
- **Key Changes**:
  - Defined `func runAutoCommit(g *engine.Git, s *engine.Summarizer) {`

[main a1b2c3d] feat(cmd): update main.go

✓ Pushed to remote
```

### 2. Generate Summary for Staged Changes (Without Committing)

Preview the commit message before committing:

```bash
git add .
git dikogit generate
```

Pass it directly to git commit:
```bash
git commit -m "$(git dikogit generate)"
```

### 3. JSON Output

For scripting and programmatic usage:

```bash
git dikogit generate --json
```

---

## Setting Up Your Custom Homebrew Tap on GitHub

To distribute `dikogit` via Homebrew:

1. **Create the Tap Repository on GitHub**:
   - Create a public repo named `homebrew-dikogit` under `Atsuwastaken` (`https://github.com/Atsuwastaken/homebrew-dikogit`).
   - Homebrew maps `brew tap Atsuwastaken/dikogit` to this repository automatically.

2. **Add the Formula**:
   - Copy `Formula/dikogit.rb` into the `Formula/` directory of your `homebrew-dikogit` repository.

3. **Tag a Release in `DIkoGit`**:
   - Create a release tag (e.g. `v2.0.0`).
   - Calculate checksum:
     ```bash
     curl -sL https://github.com/Atsuwastaken/DIkoGit/archive/refs/tags/v2.0.0.tar.gz | shasum -256
     ```
   - Update `url` and `sha256` in `Formula/dikogit.rb`.

4. **Install and Update**:
   ```bash
   brew tap Atsuwastaken/dikogit
   brew install dikogit
   ```

---

## Project Structure

```
DikoGit/
├── cmd/
│   └── git-dikogit/
│       └── main.go          # CLI entrypoint (subcommands, flags)
├── internal/
│   └── engine/
│       ├── git.go           # Git subprocess wrapper (add, commit, push, diff)
│       ├── summarizer.go    # Commit message generator
│       └── engine_test.go   # Go test suite
├── Formula/
│   └── dikogit.rb           # Homebrew formula
├── go.mod
└── README.md
```

---

## Running Tests

```bash
go test ./...
```

---

## Uninstall

```bash
brew uninstall dikogit
brew untap Atsuwastaken/dikogit
```
