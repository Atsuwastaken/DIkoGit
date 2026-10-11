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
           "-o", bin/"git-dikogit",
           "./cmd/git-dikogit"
    bin.install_symlink bin/"git-dikogit" => "dikogit"
  end

  def caveats
    <<~EOS
      Usage:
        Fully automated: stage all, commit with summary, and push:
          git dikogit

        Generate a summary for staged changes without committing:
          git dikogit generate

        JSON output for scripting:
          git dikogit generate --json
    EOS
  end

  test do
    assert_match "git-dikogit 2.0.0", shell_output("#{bin}/git-dikogit --version")
  end
end
