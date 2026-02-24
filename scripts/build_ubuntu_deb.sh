#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ARCH="${1:-amd64}"
VERSION="${VERSION:-$(git -C "$ROOT_DIR" rev-parse --short HEAD)}"
DEB_VERSION="${VERSION#v}"
DEB_VERSION="${DEB_VERSION//-/.}"
if [[ ! "$DEB_VERSION" =~ ^[0-9] ]]; then
  DEB_VERSION="0.0.0+${DEB_VERSION}"
fi
OUT_DIR="$ROOT_DIR/dist"
PKG_ROOT="$OUT_DIR/deb-build"
PKG_DIR="$PKG_ROOT/stashclip_${DEB_VERSION}_${ARCH}"
DEB_FILE="$OUT_DIR/stashclip_${DEB_VERSION}_${ARCH}.deb"
export GOCACHE="${GOCACHE:-/tmp/go-build-cache}"

case "$ARCH" in
  amd64|arm64) ;;
  *)
    echo "unsupported arch: $ARCH (use amd64 or arm64)" >&2
    exit 1
    ;;
esac

if ! command -v dpkg-deb >/dev/null 2>&1; then
  echo "dpkg-deb is required to build .deb packages" >&2
  exit 1
fi

rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR/DEBIAN" "$PKG_DIR/usr/bin" "$PKG_DIR/usr/lib/systemd/user"

CGO_ENABLED=0 GOOS=linux GOARCH="$ARCH" go build -trimpath -ldflags="-s -w" -o "$PKG_DIR/usr/bin/stashclip" "$ROOT_DIR/cmd/stashclip"
cp "$ROOT_DIR/packaging/ubuntu/stashclip.service" "$PKG_DIR/usr/lib/systemd/user/stashclip.service"

cat > "$PKG_DIR/usr/bin/stashclip-popup" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
exec /usr/bin/stashclip popup
EOF
chmod 0755 "$PKG_DIR/usr/bin/stashclip-popup"

cat > "$PKG_DIR/DEBIAN/control" <<EOF
Package: stashclip
Version: ${DEB_VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: Stashclip Team <noreply@example.com>
Depends: xclip, wl-clipboard, zenity | yad | kdialog
Description: Clipboard manager with daemon and on-demand popup picker
 Stashclip captures clipboard changes in background and provides
 an on-demand popup to select and copy saved clipboard history.
EOF

cat > "$PKG_DIR/DEBIAN/postinst" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
cat <<MSG
Stashclip instalado.
Execute no usuario que vai usar o app:
  systemctl --user daemon-reload
  systemctl --user enable --now stashclip.service

Atalho global sugerido:
  /usr/bin/stashclip-popup
MSG
EOF
chmod 0755 "$PKG_DIR/DEBIAN/postinst"

dpkg-deb --build --root-owner-group "$PKG_DIR" "$DEB_FILE"
sha256sum "$DEB_FILE" > "$DEB_FILE.sha256"

echo "deb created:"
echo "  $DEB_FILE"
echo "  $DEB_FILE.sha256"
