#!/usr/bin/env bash
set -euo pipefail

if command -v opencode >/dev/null 2>&1; then
  echo "OpenCode is already installed."
  exit 0
fi

curl -fsSL https://opencode.ai/v2/install | bash

if ! command -v opencode >/dev/null 2>&1 && [[ ! -x "$HOME/.opencode/bin/opencode" ]]; then
  echo "OpenCode installer finished, but its binary was not found on PATH."
  echo "Reopen the terminal and check the installation before starting the web UI."
  exit 1
fi

echo "OpenCode installed successfully."
