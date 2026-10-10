#!/usr/bin/env python3
import unittest
import os
import sys
import shutil
import tempfile
import subprocess
import json
from pathlib import Path

# Add bin to path so we can import git-dikogits
BIN_DIR = Path(__file__).resolve().parent.parent / "bin"
sys.path.insert(0, str(BIN_DIR))

# Import Dikogits engine
from importlib.machinery import SourceFileLoader
git_dikogits_path = str(BIN_DIR / "git-dikogits")
dikogits = SourceFileLoader("git_dikogits", git_dikogits_path).load_module()

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

    def test_hook_installation(self):
        engine = dikogits.DikogitsEngine(self.temp_dir)
        dikogits.install_hooks(is_global=False)

        git_dir = engine.get_git_dir()
        pre_commit = os.path.join(git_dir, "hooks", "pre-commit")
        prep_msg = os.path.join(git_dir, "hooks", "prepare-commit-msg")

        self.assertTrue(os.path.exists(pre_commit))
        self.assertTrue(os.path.exists(prep_msg))
        self.assertTrue(os.access(pre_commit, os.X_OK))
        self.assertTrue(os.access(prep_msg, os.X_OK))

        dikogits.uninstall_hooks(is_global=False)
        self.assertFalse(os.path.exists(pre_commit))
        self.assertFalse(os.path.exists(prep_msg))

    def test_end_to_end_commit_workflow(self):
        engine = dikogits.DikogitsEngine(self.temp_dir)
        dikogits.install_hooks(is_global=False)

        # Stage a new file
        file1 = Path(self.temp_dir) / "README.md"
        file1.write_text("# Project Docs\nThis is documentation.\n")
        subprocess.run(["git", "add", "README.md"], check=True)

        # Trigger pre-commit hook
        saved = engine.save_staged_summary()
        self.assertTrue(saved)

        dikogits_msg = os.path.join(engine.get_git_dir(), "DIKOGITS_MSG")
        self.assertTrue(os.path.exists(dikogits_msg))

        # Trigger prepare-commit-msg hook
        commit_msg_file = os.path.join(self.temp_dir, ".git", "COMMIT_EDITMSG")
        Path(commit_msg_file).write_text("# Initial commit message\n")

        applied = engine.apply_prepare_commit_msg(commit_msg_file)
        self.assertTrue(applied)

        # Verify COMMIT_EDITMSG content
        msg_content = Path(commit_msg_file).read_text()
        self.assertIn("docs", msg_content)
        self.assertIn("Dikogits Commit Summary", msg_content)
        self.assertIn("README.md", msg_content)

if __name__ == "__main__":
    unittest.main()
