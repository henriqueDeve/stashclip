#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ARCH="${1:-amd64}"
VERSION="${VERSION:-$(git -C "$ROOT_DIR" rev-parse --short HEAD)}"
OUT_DIR="$ROOT_DIR/dist"
PKG_NAME="stashclip-${VERSION}-ubuntu-${ARCH}"
STAGE_DIR="$OUT_DIR/$PKG_NAME"
export GOCACHE="${GOCACHE:-/tmp/go-build-cache}"

case "$ARCH" in
  amd64|arm64) ;;
  *)
    echo "unsupported arch: $ARCH (use amd64 or arm64)" >&2
    exit 1
    ;;
esac

rm -rf "$STAGE_DIR"
mkdir -p "$STAGE_DIR/bin"

CGO_ENABLED=0 GOOS=linux GOARCH="$ARCH" go build -trimpath -ldflags="-s -w" -o "$STAGE_DIR/bin/stashclip" "$ROOT_DIR/cmd/stashclip"

cp "$ROOT_DIR/packaging/ubuntu/install.sh" "$STAGE_DIR/install.sh"
cp "$ROOT_DIR/packaging/ubuntu/uninstall.sh" "$STAGE_DIR/uninstall.sh"
cp "$ROOT_DIR/packaging/ubuntu/stashclip.service" "$STAGE_DIR/stashclip.service"
cp "$ROOT_DIR/packaging/ubuntu/README.md" "$STAGE_DIR/README.md"
chmod +x "$STAGE_DIR/install.sh" "$STAGE_DIR/uninstall.sh"

tar -C "$OUT_DIR" -czf "$OUT_DIR/${PKG_NAME}.tar.gz" "$PKG_NAME"
sha256sum "$OUT_DIR/${PKG_NAME}.tar.gz" > "$OUT_DIR/${PKG_NAME}.tar.gz.sha256"

echo "bundle created:"
echo "  $OUT_DIR/${PKG_NAME}.tar.gz"
echo "  $OUT_DIR/${PKG_NAME}.tar.gz.sha256"
