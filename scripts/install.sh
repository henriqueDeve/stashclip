#!/usr/bin/env bash
set -euo pipefail

REPO="${STASHCLIP_REPO:-henrique/stashclip}"
VERSION="${1:-latest}"

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required" >&2
  exit 1
fi
if ! command -v sudo >/dev/null 2>&1; then
  echo "sudo is required" >&2
  exit 1
fi
if ! command -v apt-get >/dev/null 2>&1; then
  echo "This installer currently supports Ubuntu/Debian (apt-get)." >&2
  exit 1
fi

arch="$(dpkg --print-architecture)"
case "$arch" in
  amd64|arm64) ;;
  *)
    echo "unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

if [[ "$VERSION" == "latest" ]]; then
  api_url="https://api.github.com/repos/${REPO}/releases/latest"
else
  api_url="https://api.github.com/repos/${REPO}/releases/tags/${VERSION}"
fi

echo "Fetching release metadata from $api_url"
release_json="$(curl -fsSL "$api_url")"
deb_url="$(printf '%s' "$release_json" \
  | grep -Eo '"browser_download_url":[[:space:]]*"[^"]+\.deb"' \
  | cut -d'"' -f4 \
  | grep -E "_${arch}\.deb$" \
  | head -n1)"

if [[ -z "$deb_url" ]]; then
  echo "Could not find .deb for architecture ${arch} in release ${VERSION}" >&2
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
deb_file="$tmp_dir/stashclip.deb"

echo "Downloading $deb_url"
curl -fL "$deb_url" -o "$deb_file"

sudo apt-get update
sudo apt-get install -y "$deb_file"

if command -v systemctl >/dev/null 2>&1; then
  systemctl --user daemon-reload || true
  systemctl --user enable --now stashclip.service || true
fi

echo "Installed successfully."
echo "Run popup with: /usr/bin/stashclip-popup"
