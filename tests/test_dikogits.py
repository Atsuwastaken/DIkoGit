#!/usr/bin/env python3
import unittest
import os
import sys
import shutil
import tempfile
import subprocess
import json
from pathlib import Path

# Add bin to path so we can import git-dikogit
BIN_DIR = Path(__file__).resolve().parent.parent / "bin"
sys.path.insert(0, str(BIN_DIR))

# Import Dikogits engine
from importlib.machinery import SourceFileLoader
git_dikogits_path = str(BIN_DIR / "git-dikogit")
dikogits = SourceFileLoader("git_dikogit", git_dikogits_path).load_module()

class TestDikogits(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.mkdtemp()
        self.old_cwd = os.getcwd()
        os.chdir(self.temp_dir)
        
        # Init git repo
        subprocess.run(["git", "init"], check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        subprocess.run(["git", "config", "user.name", "Test User"], check=True)
        subprocess.run(["git", "config", "user.email", "test@example.com"], check=True)

    def tearDown(self):
        os.chdir(self.old_cwd)
        shutil.rmtree(self.temp_dir)

    def test_parse_staged_changes(self):
        engine = dikogits.DikogitsEngine(self.temp_dir)

        # Create sample files
        file1 = Path(self.temp_dir) / "app.py"
        file1.write_text("def main():\n    print('Hello World')\n")
        
        subprocess.run(["git", "add", "app.py"], check=True)

        changes = engine.parse_changes()
        self.assertEqual(len(changes['files']), 1)
        self.assertEqual(changes['files'][0]['path'], "app.py")
        self.assertEqual(changes['files'][0]['status'], "added")
        self.assertGreater(changes['additions'], 0)

    def test_generate_summary_classification(self):
        engine = dikogits.DikogitsEngine(self.temp_dir)

        # 1. Feature commit
        file1 = Path(self.temp_dir) / "feature.py"
        file1.write_text("def new_feature():\n    pass\n")
        subprocess.run(["git", "add", "feature.py"], check=True)

        summary = engine.generate_summary()
        self.assertIsNotNone(summary)
        self.assertTrue(summary['title'].startswith("feat"))
        self.assertIn("feature.py", summary['body'])

    def test_cli_execution_generate(self):
        # Stage a documentation file
        readme = Path(self.temp_dir) / "README.md"
        readme.write_text("# Project Documentation\nTesting CLI generation.\n")
        subprocess.run(["git", "add", "README.md"], check=True)

        # Run git-dikogit CLI generate command
        result = subprocess.run(
            [git_dikogits_path, "generate"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=True
        )
        self.assertIn("docs", result.stdout)
        self.assertIn("README.md", result.stdout)
        self.assertIn("Dikogits Commit Summary", result.stdout)

    def test_cli_generate_no_staged_changes(self):
        result = subprocess.run(
            [git_dikogits_path, "generate"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("No staged changes found", result.stderr)

    def test_stage_all(self):
        engine = dikogits.DikogitsEngine(self.temp_dir)
        file1 = Path(self.temp_dir) / "test.txt"
        file1.write_text("sample content")

        success = engine.stage_all()
        self.assertTrue(success)

        # Check git status shows staged file
        res = subprocess.run(["git", "diff", "--cached", "--name-only"], capture_output=True, text=True, check=True)
        self.assertIn("test.txt", res.stdout)

    def test_auto_commit_engine(self):
        engine = dikogits.DikogitsEngine(self.temp_dir)
        file1 = Path(self.temp_dir) / "app.py"
        file1.write_text("def hello():\n    return 'world'\n")

        res = engine.auto_commit()
        self.assertEqual(res["status"], "success")
        self.assertIn("feat", res["summary"]["title"])

        # Check git log has the commit
        log_res = subprocess.run(["git", "log", "-1", "--pretty=%B"], capture_output=True, text=True, check=True)
        self.assertIn(res["summary"]["title"], log_res.stdout)
        self.assertIn("Dikogits Commit Summary", log_res.stdout)

    def test_cli_execution_auto_commit_default(self):
        # Create an unstaged change (not using git add beforehand)
        code_file = Path(self.temp_dir) / "calculator.py"
        code_file.write_text("def add(a, b):\n    return a + b\n")

        # Running `git-dikogit` with no args should automatically `git add .` and commit
        result = subprocess.run(
            [git_dikogits_path],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=True
        )
        self.assertIn("Staged all changes", result.stdout)
        self.assertIn("Committed with summary", result.stdout)
        self.assertIn("calculator.py", result.stdout)

        # Verify git log
        log_res = subprocess.run(["git", "log", "-1", "--pretty=%B"], capture_output=True, text=True, check=True)
        self.assertIn("calculator.py", log_res.stdout)

    def test_cli_execution_auto_commit_explicit_command(self):
        code_file = Path(self.temp_dir) / "utils.py"
        code_file.write_text("def format_str():\n    pass\n")

        result = subprocess.run(
            [git_dikogits_path, "commit"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=True
        )
        self.assertIn("Staged all changes", result.stdout)
        self.assertIn("Committed with summary", result.stdout)

    def test_cli_auto_commit_no_changes(self):
        # No files created or modified
        result = subprocess.run(
            [git_dikogits_path],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("No changes found to commit (working tree clean)", result.stderr)

    def test_cli_not_a_git_repo(self):
        non_git_dir = tempfile.mkdtemp()
        try:
            result = subprocess.run(
                [git_dikogits_path],
                cwd=non_git_dir,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("fatal: not a git repository", result.stderr)
        finally:
            shutil.rmtree(non_git_dir)

if __name__ == "__main__":
    unittest.main()
