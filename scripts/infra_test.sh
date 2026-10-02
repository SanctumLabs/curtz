#!/usr/bin/env bash
# Tests for scripts/infra.sh and scripts/infra_env_check.sh. Needs no Docker: infra.sh runs with DRY_RUN=1.
# shellcheck disable=SC2016  # the fixtures below write a literal ${VAR:-default} and $ on purpose
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1

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

# --- infra.sh wait (a stub `docker` supplies canned compose output) ---------------------------------------------
mkdir -p "$tmp/bin"
cat >"$tmp/bin/docker" <<'STUB'
#!/usr/bin/env bash
case "$*" in
  *"config --services"*) printf '%s\n' kafka-init-single kafka-single kafka-ui ;;
  *" ps "*) cat "$STUB_STATUS" ;;
esac
STUB
chmod +x "$tmp/bin/docker"

wait_with() { # one status line per service: service|state|health|exitcode
  printf '%s\n' "$1" >"$tmp/status"
  env -u DRY_RUN PATH="$tmp/bin:$PATH" STUB_STATUS="$tmp/status" WAIT_TIMEOUT=0 scripts/infra.sh wait kafka single
}

job_running='kafka-single|running|healthy|0
kafka-init-single|running||0
kafka-ui|running||0'
job_done='kafka-single|running|healthy|0
kafka-init-single|exited||0
kafka-ui|running||0'
job_failed='kafka-single|running|healthy|0
kafka-init-single|exited||1
kafka-ui|running||0'
broker_starting='kafka-single|running|starting|0
kafka-init-single|exited||0
kafka-ui|running||0'

expect_exit "wait: a still-running one-shot job is not ready" 1 wait_with "$job_running"
out="$(wait_with "$job_running" 2>&1)"
case "$out" in *"waiting for: kafka-init-single"*) pass "wait: the message names the pending job" ;; *) fail "wait: the message names the pending job ($out)" ;; esac
expect_exit "wait: ready when the job exited 0 and a service without a healthcheck is running" 0 wait_with "$job_done"
expect_exit "wait: a job that exited non-zero is not ready" 1 wait_with "$job_failed"
expect_exit "wait: a service still starting is not ready" 1 wait_with "$broker_starting"

if [ "$failures" -ne 0 ]; then
  echo "$failures failure(s)"
  exit 1
fi
echo "all tests passed"
