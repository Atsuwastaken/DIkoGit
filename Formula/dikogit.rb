class Dikogit < Formula
  desc "Automated commit and summary generator CLI extension for Git on macOS"
  homepage "https://github.com/Atsuwastaken/DIkoGit"
  url "https://github.com/Atsuwastaken/DIkoGit/archive/refs/tags/v1.1.0.tar.gz"
  # sha256 "<INSERT_RELEASE_SHA256_HERE>"
  license "MIT"
  head "https://github.com/Atsuwastaken/DIkoGit.git", branch: "main"

  uses_from_macos "python"

  def install
    bin.install "bin/git-dikogit"
    bin.install_symlink bin/"git-dikogit" => "dikogit"
  end

  def caveats
    <<~EOS
      Usage:
        To automatically stage all changes and commit with summary:
          git dikogit

        To only generate summary for staged changes without committing:
          git dikogit generate
    EOS
  end

  test do
    assert_match "git-dikogit 1.1.0", shell_output("#{bin}/git-dikogit --version")
  end
end
