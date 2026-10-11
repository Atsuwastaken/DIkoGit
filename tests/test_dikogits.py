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

    def test_cli_no_staged_changes(self):
        result = subprocess.run(
            [git_dikogits_path, "generate"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("No staged changes found", result.stderr)

if __name__ == "__main__":
    unittest.main()
