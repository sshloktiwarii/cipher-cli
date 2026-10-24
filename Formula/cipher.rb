class Cipher < Formula
  desc "Zero-dependency, edge-native autonomous AI PaaS"
  homepage "https://github.com/sshloktiwarii/cipher-cli"
  version "1.2.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/sshloktiwarii/cipher-cli/releases/download/v1.2.0/cipher_darwin_arm64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    else
      url "https://github.com/sshloktiwarii/cipher-cli/releases/download/v1.2.0/cipher_darwin_amd64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/sshloktiwarii/cipher-cli/releases/download/v1.2.0/cipher_linux_arm64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    else
      url "https://github.com/sshloktiwarii/cipher-cli/releases/download/v1.2.0/cipher_linux_amd64.tar.gz"
      sha256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    end
  end

  def install
    bin.install "cipher"
  end

  def post_install
    (var/"cipher").mkpath
    chmod 0700, var/"cipher"
  end

  test do
    system "#{bin}/cipher", "--help"
  end
end
