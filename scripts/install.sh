#!/bin/bash
set -euo pipefail

REPO="bucketeer-io/code-refs"
BINARY="bucketeer-find-code-refs"

# Detect OS
OS=$(uname -s)
if [ "$OS" != "Linux" ]; then
    echo "This script supports Linux only."
    echo "For macOS, use Homebrew:"
    echo "  brew tap bucketeer-io/code-refs https://github.com/bucketeer-io/code-refs"
    echo "  brew install bucketeer-find-code-refs"
    exit 1
fi

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)        PKG_ARCH="amd64"; TAR_ARCH="amd64" ;;
    aarch64|arm64) PKG_ARCH="arm64"; TAR_ARCH="arm64" ;;
    i386|i686)     PKG_ARCH="i386";  TAR_ARCH="386"   ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get latest version from GitHub
VERSION=$(curl -sf "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | cut -d'"' -f4)
VERSION_NUM="${VERSION#v}"

echo "Installing $BINARY $VERSION for Linux/$ARCH..."

BASE_URL="https://github.com/$REPO/releases/download/$VERSION"

install_deb() {
    TMP=$(mktemp /tmp/bucketeer-code-refs.XXXXXX.deb)
    curl -sfL "$BASE_URL/code-refs_${VERSION_NUM}.${PKG_ARCH}.deb" -o "$TMP"
    dpkg -i "$TMP"
    rm -f "$TMP"
}

install_rpm() {
    TMP=$(mktemp /tmp/bucketeer-code-refs.XXXXXX.rpm)
    curl -sfL "$BASE_URL/code-refs_${VERSION_NUM}.${PKG_ARCH}.rpm" -o "$TMP"
    rpm -i "$TMP"
    rm -f "$TMP"
}

install_tar() {
    TMP_DIR=$(mktemp -d)
    curl -sfL "$BASE_URL/code-refs_${VERSION_NUM}_linux_${TAR_ARCH}.tar.gz" \
        | tar -xz -C "$TMP_DIR"
    install -m 755 "$TMP_DIR/$BINARY" /usr/local/bin/
    rm -rf "$TMP_DIR"
}

if command -v dpkg &>/dev/null; then
    install_deb
elif command -v rpm &>/dev/null; then
    install_rpm
else
    install_tar
fi

echo ""
echo "$BINARY $VERSION installed successfully!"