# Dikogit - Git CLI Extension for macOS

**Dikogit** is a Git command-line extension for macOS that automatically analyzes staged code changes, categorizes the commit type, and generates structured commit summaries.

---

## Features

- 🏷️ **Conventional Commit Categorization**: Automatically categorizes commits into standard types (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`) based on changed files and diff inspection.
- 📊 **Detailed Statistics & Breakdown**: Computes file additions, deletions, status breakdown (`[A]`, `[M]`, `[D]`, `[R]`), and identifies key constructs (functions, imports, section titles).
- 🍺 **Homebrew Distribution**: Easy global installation and version management via Homebrew tap.
- 💻 **Simple CLI Interface**: Use `git dikogit generate` or pipe it directly into `git commit`.
- ⚡ **Zero External Dependencies**: Pure Python 3 using standard library modules.

---

## Installation via Homebrew (Recommended)

Install `dikogit` globally on macOS:

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

### 1. Generate Summary for Staged Changes
Stage your changes with `git add`, then run:

```bash
git add .
git dikogit generate
```

Example Output:
```text
feat(root): update 2 files (+25 -2)

### Dikogits Commit Summary
- **Stats**: 2 file(s) changed, 25 insertion(s)(+), 2 deletion(s)(-)
- **File Breakdown**:
  - [M] `file1.py` (+20 -2)
  - [A] `docs/README.md` (+5 -0)
- **Key Changes**:
  - Defined `def calculate_summary():`
```

### 2. Commit Directly with Generated Summary
You can directly pass the generated summary to `git commit`:

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
   - Create a release tag (e.g. `v1.0.0`).
   - Calculate checksum:
     ```bash
     curl -sL https://github.com/Atsuwastaken/DIkoGit/archive/refs/tags/v1.0.0.tar.gz | shasum -256
     ```
   - Update `url` and `sha256` in `Formula/dikogit.rb`.

4. **Install and Update**:
   ```bash
   brew tap Atsuwastaken/dikogit
   brew install dikogit
   ```

---

## Local Development Installation

If you prefer installing without Homebrew during local development:

```bash
./install.sh
```

Ensure `~/.local/bin` is in your `$PATH`:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

To remove the local symlink:
```bash
./uninstall.sh
```

---

## Running Tests

Run the test suite:

```bash
python3 -m unittest discover -s tests
```
