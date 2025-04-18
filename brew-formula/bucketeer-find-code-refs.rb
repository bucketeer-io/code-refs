# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class BucketeerFindCodeRefs < Formula
  desc "Bucketeer Code References. This tool will find and send feature flag code references to Bucketeer's backend."
  homepage "https://bucketeer.io"
  version "0.0.2"
  license "Apache 2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/bucketeer-io/code-refs/releases/download/v0.0.2/code-refs_0.0.2_darwin_amd64.tar.gz"
      sha256 "535880062c53758c8366da58114ea9adc17f902b87b0e0471df670c0520dd9c4"

      def install
        bin.install "bucketeer-find-code-refs"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/bucketeer-io/code-refs/releases/download/v0.0.2/code-refs_0.0.2_darwin_arm64.tar.gz"
      sha256 "e9077a39dfc07908e9af3bf09246498b9618138f1ac60873e48d9eab905059d8"

      def install
        bin.install "bucketeer-find-code-refs"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/bucketeer-io/code-refs/releases/download/v0.0.2/code-refs_0.0.2_linux_amd64.tar.gz"
        sha256 "53340950e9c37dd45bc783993cb7fd79c3616e78e69bace721598eb3e2da534e"

        def install
          bin.install "bucketeer-find-code-refs"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/bucketeer-io/code-refs/releases/download/v0.0.2/code-refs_0.0.2_linux_arm64.tar.gz"
        sha256 "dd8230aa38ae24af17c6cdb3b4015cbf096e1f23b7e3b1b8e92b5d4759eab4b7"

        def install
          bin.install "bucketeer-find-code-refs"
        end
      end
    end
  end
end
