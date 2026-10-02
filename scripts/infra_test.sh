#!/usr/bin/env bash
# Tests for scripts/infra.sh and scripts/infra_env_check.sh. Needs no Docker: infra.sh runs with DRY_RUN=1.
set -uo pipefail
cd "$(dirname "$0")/.."

failures=0
pass() { echo "ok:   $1"; }
fail() { echo "FAIL: $1"; failures=$((failures + 1)); }

expect_output() { # name expected command...
  local name="$1" expected="$2"
  shift 2
  local actual
  actual="$("$@" 2>&1)"
  if [ "$actual" = "$expected" ]; then
    pass "$name"
  else
    fail "$name"
    printf -- '--- expected\n%s\n--- actual\n%s\n' "$expected" "$actual"
  fi
}

expect_exit() { # name expected-code command...
  local name="$1" want="$2"
  shift 2
  "$@" >/dev/null 2>&1
  local got=$?
  if [ "$got" -eq "$want" ]; then pass "$name"; else fail "$name (exit $got, want $want)"; fi
}

export DRY_RUN=1 COMPOSE="docker compose"

# --- infra.sh ---------------------------------------------------------------------------------------------------
expect_output "up kafka single removes the ha containers first" \
"+ docker compose --profile kafka-ha rm -sf
+ docker compose --profile kafka-single up -d
+ wait kafka-single" \
  scripts/infra.sh up kafka single

expect_output "mode defaults to ha" \
"+ docker compose --profile kafka-single rm -sf
+ docker compose --profile kafka-ha up -d
+ wait kafka-ha" \
  scripts/infra.sh up kafka

expect_output "aggregate stacks follow the mode" \
"+ docker compose --profile core-ha rm -sf
+ docker compose --profile core-single up -d
+ wait core-single" \
  scripts/infra.sh up core single

expect_output "observability has no mode and nothing to remove" \
"+ docker compose --profile observability up -d
+ wait observability" \
  scripts/infra.sh up observability single

expect_output "down removes both modes of a stack" \
"+ docker compose --profile redis-ha --profile redis-single rm -sf" \
  scripts/infra.sh down redis

expect_output "down legacy" \
"+ docker compose --profile legacy rm -sf" \
  scripts/infra.sh down legacy

expect_output "wait resolves the profile" \
"+ wait postgres-single" \
  scripts/infra.sh wait postgres single

expect_exit "unknown stack is a usage error" 2 scripts/infra.sh up nope
expect_exit "unknown mode is a usage error" 2 scripts/infra.sh up kafka triple
expect_exit "missing arguments is a usage error" 2 scripts/infra.sh
expect_exit "unknown action is a usage error" 2 scripts/infra.sh restart kafka

# --- infra_env_check.sh -----------------------------------------------------------------------------------------
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
printf '# --- Local infrastructure\nPG_APP_PASSWORD=it%ss\n' "'" >"$tmp/unsafe-quote.env"
printf '# --- Local infrastructure\nPG_APP_PASSWORD=pa$ss\n' >"$tmp/unsafe-dollar.env"
printf '# --- Local infrastructure\nPG_APP_PASSWORD=safe\n' >"$tmp/safe.env"
printf 'x: ${PG_APP_PASSWORD:-safe}\n' >"$tmp/match.yml"
printf 'x: ${PG_APP_PASSWORD:-other}\n' >"$tmp/mismatch.yml"
printf 'x: ${PG_UNKNOWN:-v}\n' >"$tmp/unknown.yml"

expect_exit "env check rejects a quote in a value" 1 \
  env ENV_EXAMPLE="$tmp/unsafe-quote.env" COMPOSE_FILES=/dev/null scripts/infra_env_check.sh
expect_exit "env check rejects a dollar sign in a value" 1 \
  env ENV_EXAMPLE="$tmp/unsafe-dollar.env" COMPOSE_FILES=/dev/null scripts/infra_env_check.sh
expect_exit "env check accepts safe values with matching defaults" 0 \
  env ENV_EXAMPLE="$tmp/safe.env" COMPOSE_FILES="$tmp/match.yml" scripts/infra_env_check.sh
expect_exit "env check rejects a default that differs from .env.example" 1 \
  env ENV_EXAMPLE="$tmp/safe.env" COMPOSE_FILES="$tmp/mismatch.yml" scripts/infra_env_check.sh
expect_exit "env check rejects a variable missing from .env.example" 1 \
  env ENV_EXAMPLE="$tmp/safe.env" COMPOSE_FILES="$tmp/unknown.yml" scripts/infra_env_check.sh
expect_exit "env check passes on the repository" 0 scripts/infra_env_check.sh

if [ "$failures" -ne 0 ]; then
  echo "$failures failure(s)"
  exit 1
fi
echo "all tests passed"
