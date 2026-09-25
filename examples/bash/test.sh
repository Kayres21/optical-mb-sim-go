#!/bin/bash

SESSION="make_cluster"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CONFIG_DIR="$REPO_ROOT/configs/TEST"

mapfile -t CONFIGS < <(find "$CONFIG_DIR" -type f -name '*.json' | sort)

if [ "${#CONFIGS[@]}" -eq 0 ]; then
    echo "No config files found under $CONFIG_DIR" >&2
    exit 1
fi

for i in "${!CONFIGS[@]}"; do
    CONFIGS[$i]="${CONFIGS[$i]#$REPO_ROOT/}"
done

# 1. Start the first window
first_config="${CONFIGS[0]}"
tmux new-session -d -s "$SESSION" -n "win-1" "cd '$REPO_ROOT' && make run CONFIG=$first_config; exec \$SHELL"
sleep 5

# 2. Create one window per remaining config file
for i in $(seq 1 $(( ${#CONFIGS[@]} - 1 ))); do
    win_num=$((i + 1))
    config_path="${CONFIGS[$i]}"
    tmux new-window -t "$SESSION" -n "win-$win_num" "cd '$REPO_ROOT' && make run CONFIG=$config_path; exec \$SHELL"
    sleep 5
done

echo "Started ${#CONFIGS[@]} config-driven make tasks in session '$SESSION'."