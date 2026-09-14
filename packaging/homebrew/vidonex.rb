class Vidonex < Formula
  desc "Declarative Video Composition & FFmpeg Filtergraph Compiler Engine in Go"
  homepage "https://farshidrezaei.github.io/vidonex/"
  version "1.0.0"
  license "MIT"

  on_macos do
    url "https://github.com/farshidrezaei/vidonex/releases/download/v1.0.0/vidonex-cli-darwin-universal"
    sha256 "REPLACE_WITH_DARWIN_SHA256"

    def install
      bin.install "vidonex-cli-darwin-universal" => "vidonex"
    end
  end

  on_linux do
    url "https://github.com/farshidrezaei/vidonex/releases/download/v1.0.0/vidonex-cli-linux-amd64"
    sha256 "REPLACE_WITH_LINUX_SHA256"

    def install
      bin.install "vidonex-cli-linux-amd64" => "vidonex"
    end
  end

  test do
    assert_match "Vidonex", shell_output("#{bin}/vidonex --help")
  end
end
