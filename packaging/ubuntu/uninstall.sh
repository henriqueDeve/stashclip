#!/usr/bin/env bash
set -euo pipefail

BIN_DST="$HOME/.local/bin/stashclip"
POPUP_WRAPPER="$HOME/.local/bin/stashclip-popup"
SERVICE_DST="$HOME/.config/systemd/user/stashclip.service"

if command -v systemctl >/dev/null 2>&1; then
  systemctl --user disable --now stashclip.service 2>/dev/null || true
  systemctl --user daemon-reload || true
fi

rm -f "$SERVICE_DST" "$BIN_DST" "$POPUP_WRAPPER"
echo "stashclip removido"
