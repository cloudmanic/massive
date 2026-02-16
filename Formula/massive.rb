#
# Date: 2026-02-15
# Copyright (c) 2026. All rights reserved.
#

# Homebrew formula for the Massive CLI. Downloads pre-built binaries
# from GitHub releases for the current platform and architecture.
class Massive < Formula
  desc "CLI for the Massive financial data API"
  homepage "https://github.com/cloudmanic/massive"
  license "MIT"
  version "latest"

  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/massive/releases/latest/download/massive-darwin-arm64"
  elsif OS.mac? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/massive/releases/latest/download/massive-darwin-amd64"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/cloudmanic/massive/releases/latest/download/massive-linux-arm64"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/cloudmanic/massive/releases/latest/download/massive-linux-amd64"
  end

  def install
    binary_name = stable.url.split("/").last
    bin.install binary_name => "massive"
  end

  test do
    assert_match "massive version", shell_output("#{bin}/massive --version")
  end
end
