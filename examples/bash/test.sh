#!/bin/bash

SESSION="make_cluster"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

mapfile -t CONFIGS < <(find "$REPO_ROOT/configs" -type f -name '*.json' | sort)

if [ "${#CONFIGS[@]}" -eq 0 ]; then
    echo "No config files found under $REPO_ROOT/configs" >&2
    exit 1
fi

for i in "${!CONFIGS[@]}"; do
    CONFIGS[$i]="${CONFIGS[$i]#$REPO_ROOT/}"
done

# 1. Start the first window
first_config="${CONFIGS[0]}"
tmux new-session -d -s "$SESSION" -n "win-1" "cd '$REPO_ROOT' && make run CONFIG=$first_config; exec \$SHELL"

# 2. Create one window per remaining config file
for i in $(seq 1 $(( ${#CONFIGS[@]} - 1 ))); do
    win_num=$((i + 1))
    config_path="${CONFIGS[$i]}"
    tmux new-window -t "$SESSION" -n "win-$win_num" "cd '$REPO_ROOT' && make run CONFIG=$config_path; exec \$SHELL"
done

echo "Started ${#CONFIGS[@]} config-driven make tasks in session '$SESSION'."