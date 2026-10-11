class Dikogit < Formula
  desc "Fully automated Git CLI extension: stage, commit, and push with one command"
  homepage "https://github.com/Atsuwastaken/DIkoGit"
  url "https://github.com/Atsuwastaken/DIkoGit/archive/refs/tags/v2.0.0.tar.gz"
  # sha256 "<INSERT_RELEASE_SHA256_HERE>"
  license "MIT"
  head "https://github.com/Atsuwastaken/DIkoGit.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build",
           "-ldflags", "-s -w",
           "-o", bin/"git-dikotgit",
           "./cmd/git-dikogit"
    bin.install_symlink bin/"git-dikotgit" => "dikotgit"
  end

  def caveats
    <<~EOS
      Usage:
        Fully automated: stage all, commit with summary, and push:
          git dikotgit

        Generate a summary for staged changes without committing:
          git dikotgit generate

        JSON output for scripting:
          git dikotgit generate --json
    EOS
  end

  test do
    assert_match "git-dikotgit 2.0.0", shell_output("#{bin}/git-dikotgit --version")
  end
end
