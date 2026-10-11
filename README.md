# Dikogit - Git CLI Extension for macOS

**Dikogit** is a Git command-line extension for macOS that automatically analyzes changed files, categorizes the commit type, and generates structured commit summaries.

---

## Features

- 🚀 **Fully Automated Commits**: Run `git dikogit` to automatically stage all changes (`git add .`) and commit directly with a detailed summary.
- 🏷️ **Conventional Commit Categorization**: Automatically categorizes commits into standard types (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`) based on changed files and diff inspection.
- 📊 **Detailed Statistics & Breakdown**: Computes file additions, deletions, status breakdown (`[A]`, `[M]`, `[D]`, `[R]`), and identifies key constructs (functions, imports, section titles).
- 🍺 **Homebrew Distribution**: Easy global installation and version management via Homebrew tap.
- 💻 **Flexible CLI Interface**: Auto-commit with `git dikogit` or preview/generate without committing via `git dikogit generate`.
- ⚡ **Zero External Dependencies**: Pure Python 3 using standard library modules.

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

### 1. Fully Automated Commit (Default)
To automatically stage all changed files (`git add .`), analyze the changes, and commit directly to your repository with a structured summary:

```bash
git dikogit
```

Example Output:
```text
✓ Staged all changes (git add .)
✓ Committed with summary:

feat(root): update 2 files (+25 -2)

### Dikogits Commit Summary
- **Stats**: 2 file(s) changed, 25 insertion(s)(+), 2 deletion(s)(-)
- **File Breakdown**:
  - [M] `file1.py` (+20 -2)
  - [A] `docs/README.md` (+5 -0)
- **Key Changes**:
  - Defined `def calculate_summary():`

[main 8f3c1d2] feat(root): update 2 files (+25 -2)
 2 files changed, 25 insertions(+), 2 deletions(-)
```

### 2. Generate Summary for Staged Changes (Without Committing)
If you prefer to stage changes manually and only generate the summary without committing:

```bash
git add .
git dikogit generate
```

You can also pass the generated summary directly to `git commit`:
```bash
git commit -m "$(git dikogit generate)"
```

### 3. JSON Output
For programmatic usage and scripts:

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
   - Copy [`Formula/dikogit.rb`](file:///Users/atsu/DikoGit/Formula/dikogit.rb) into the `Formula/` directory of your `homebrew-dikogit` repository.

3. **Tag a Release in `DIkoGit`**:
   - Create a release tag (e.g. `v1.1.0`).
   - Calculate checksum:
     ```bash
     curl -sL https://github.com/Atsuwastaken/DIkoGit/archive/refs/tags/v1.1.0.tar.gz | shasum -256
     ```
   - Update `url` and `sha256` in `Formula/dikogit.rb`.

4. **Install and Update**:
   ```bash
   brew tap Atsuwastaken/dikogit
   brew install dikogit
   ```

---

## Running Tests

Run the test suite:

```bash
python3 -m unittest discover -s tests
```
