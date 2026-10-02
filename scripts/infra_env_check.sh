#!/usr/bin/env bash
# Guards the "works with an empty or stale .env" promise:
#  1. every ${VAR:-default} in the compose files must use exactly the default listed in .env.example;
#  2. values in the infrastructure section of .env.example must be safe to interpolate into shell and SQL
#     (no quotes, backslashes or dollar signs).
# Env: ENV_EXAMPLE (default .env.example), COMPOSE_FILES (default: root and deploy/*/compose.yml).
set -euo pipefail

cd "$(dirname "$0")/.."
ENV_EXAMPLE="${ENV_EXAMPLE:-.env.example}"
shopt -s nullglob
default_files=(docker-compose.yml deploy/*/compose.yml)
COMPOSE_FILES="${COMPOSE_FILES:-${default_files[*]}}"
status=0

awk '
  /^# --- Local infrastructure/ { on = 1; next }
  on && /^[A-Z0-9_]+=/ {
    key = substr($0, 1, index($0, "=") - 1)
    val = substr($0, index($0, "=") + 1)
    if (index(val, "\"") || index(val, "\\") || index(val, "$") || index(val, sprintf("%c", 39))) {
      print "unsafe character in " key
      bad = 1
    }
  }
  END { exit bad }
' "$ENV_EXAMPLE" || status=1

# shellcheck disable=SC2086
while IFS='|' read -r file var def; do
  [ -n "$var" ] || continue
  want="$(grep -E "^${var}=" "$ENV_EXAMPLE" | head -1 | cut -d= -f2- || true)"
  if [ -z "$want" ]; then
    echo "$file: \${$var:-...} has no entry in $ENV_EXAMPLE"
    status=1
  elif [ "$want" != "$def" ]; then
    echo "$file: \${$var:-$def} differs from $ENV_EXAMPLE ($want)"
    status=1
  fi
done < <(grep -o -H -E '\$\{[A-Z0-9_]+:-[^}]*\}' $COMPOSE_FILES 2>/dev/null |
  sed -E 's/^([^:]+):\$\{([A-Z0-9_]+):-([^}]*)\}$/\1|\2|\3/')

exit "$status"
