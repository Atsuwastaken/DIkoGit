# Dikogits - Git Extension for macOS

**Dikogits** is a Git extension and hook integration for macOS that automatically analyzes staged code changes before each commit, generates a structured commit summary, and updates the commit message accordingly.

---

## Key Features

- ⚡ **Automatic Hook Execution**: Integrates with Git `pre-commit` and `prepare-commit-msg` hooks to capture staged diffs seamlessly.
- 🏷️ **Conventional Commit Categorization**: Automatically categorizes commits into standard types (`feat`, `fix`, `docs`, `refactor`, `test`, `chore`) based on changed files and diff inspection.
- 📊 **Detailed Breakdown**: Includes file addition/deletion statistics, file status breakdown (`[A]`, `[M]`, `[D]`, `[R]`), and key constructs added (functions, imports, sections).
- 🛠️ **Seamless CLI Integration**: Can be run manually as `git dikogit generate` or automatically upon executing `git commit`.
- 🍏 **Native macOS Compatibility**: Designed for macOS zsh/bash environments without requiring external binary dependencies.

---

## Installation

### 1. One-Step Installer Script
Run the installer script from the project directory:

```bash
./install.sh
```

This will link `git-dikogit` into `~/.local/bin` and install local git hooks for the current repository.

### 2. Manual Installation
Add `bin/git-dikogit` to your system `$PATH` or symlink it:

```bash
ln -sf $(pwd)/bin/git-dikogit ~/.local/bin/git-dikogit
```

Then install local Git hooks in any repository:

```bash
git dikogit install
```

Or install globally for all Git repositories on your macOS system:

```bash
git dikogit install --global
```

---

## Usage Workflow

1. **Stage your changes as usual**:
   ```bash
   git add file1.py docs/README.md
   ```

2. **Run `git commit`**:
   The `pre-commit` hook automatically analyzes staged changes, and `prepare-commit-msg` updates the commit message editor with the pre-commit generated summary:

   ```
   feat(root): update 2 files (+25 -2)

   ### Dikogits Commit Summary
   - **Stats**: 2 file(s) changed, 25 insertion(s)(+), 2 deletion(s)(-)
   - **File Breakdown**:
     - [M] `file1.py` (+20 -2)
     - [A] `docs/README.md` (+5 -0)
   - **Key Changes**:
     - Defined `def calculate_summary():`
   ```

3. **Manual Summary Generation**:
   You can also generate a summary anytime without committing:

   ```bash
   git dikogit generate
   ```

---

## Uninstallation

To remove hooks from a repository or globally:

```bash
./uninstall.sh
```

or

```bash
git dikogit uninstall [--global]
```

---

## Running Tests

Run the included unit test suite:

```bash
python3 -m unittest discover -s tests
```
