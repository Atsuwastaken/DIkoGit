class Dikogit < Formula
  desc "Automated commit summary generator CLI extension for Git on macOS"
  homepage "https://github.com/Atsuwastaken/DIkoGit"
  url "https://github.com/Atsuwastaken/DIkoGit/archive/refs/tags/v1.0.0.tar.gz"
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
        Stage your changes with `git add`, then run:
          git dikogit generate

        To commit directly with the generated summary:
          git commit -m "$(git dikogit generate)"
    EOS
  end

  test do
    assert_match "git-dikogit 1.0.0", shell_output("#{bin}/git-dikogit --version")
  end
end
