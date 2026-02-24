#!/usr/bin/env bash
set -euo pipefail

INSTALL_DEPS=0
for arg in "$@"; do
  case "$arg" in
    --install-deps) INSTALL_DEPS=1 ;;
    *)
      echo "unknown argument: $arg" >&2
      echo "usage: ./install.sh [--install-deps]" >&2
      exit 1
      ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_SRC="$SCRIPT_DIR/bin/stashclip"
BIN_DST="$HOME/.local/bin/stashclip"
SERVICE_DST_DIR="$HOME/.config/systemd/user"
SERVICE_DST="$SERVICE_DST_DIR/stashclip.service"
POPUP_WRAPPER="$HOME/.local/bin/stashclip-popup"

if [[ ! -x "$BIN_SRC" ]]; then
  echo "binary not found at $BIN_SRC" >&2
  exit 1
fi

if [[ "$INSTALL_DEPS" -eq 1 ]]; then
  if command -v apt-get >/dev/null 2>&1; then
    sudo apt-get update
    sudo apt-get install -y xclip wl-clipboard zenity yad kdialog
  else
    echo "apt-get not found; install dependencies manually:" >&2
    echo "  xclip wl-clipboard zenity yad kdialog" >&2
  fi
fi

mkdir -p "$HOME/.local/bin" "$SERVICE_DST_DIR"
install -m 0755 "$BIN_SRC" "$BIN_DST"
install -m 0644 "$SCRIPT_DIR/stashclip.service" "$SERVICE_DST"

cat > "$POPUP_WRAPPER" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
exec "$HOME/.local/bin/stashclip" popup
EOF
chmod +x "$POPUP_WRAPPER"

if command -v systemctl >/dev/null 2>&1; then
  systemctl --user daemon-reload
  systemctl --user enable --now stashclip.service
else
  "$BIN_DST" daemon start || true
fi

echo "stashclip instalado em $BIN_DST"
echo "daemon iniciado"
echo "atalho global sugerido:"
echo "  comando: $HOME/.local/bin/stashclip-popup"
