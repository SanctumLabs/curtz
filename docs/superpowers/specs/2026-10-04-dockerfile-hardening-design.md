# Production image, `app` compose stack and CI — design

Status: draft for review · Date: 2026-10-04 · Slice 2 of 5

## 1. Intent

The `Dockerfile` cannot build the app today and is not safe to ship. It pins Go 1.18.1 while `go.mod` needs 1.26, runs as
root, has no `HEALTHCHECK`, builds with `-mod=mod` (a build may change dependencies), and is fed by an empty `.dockerignore`,
so `.git`, `.env` and `.superpowers/` reach the builder stage. CI is out of date in the same way: `docker.yml` and `build_app.yml`
use v1/v2 actions and Go 1.18, and the workflows that run on `workflow_run` check out the default branch instead of the commit
that triggered them.

This slice produces a hardened, reproducible image of the API and the migrator, an `app` compose stack so the container runs on
the `curtz` network the way it will in production, repaired lint and scan targets, and modernized CI.

**Success criteria**

1. `docker build` produces an image that runs as uid 65532, has no shell, and carries the API, the migrator and the migrations.
2. `hadolint` is clean, and Trivy reports no fixable HIGH or CRITICAL finding in the image.
3. `make infra.app.up MODE=single|ha` starts Postgres, Redis and the app and reaches `ready`; the container is `healthy`
   and `GET /health/ready` is `ok`.
4. The container runs with a read-only root filesystem, all capabilities dropped and `no-new-privileges`.
5. `docker run --entrypoint /app/migrator` applies the migrations from the image (a second run is "no change").
6. The build context contains no `.env`, `.git` or test files.
7. Every workflow passes `actionlint`, and the Docker workflow scans before it publishes.

## 2. Scope

Slice table (updated; slice 3 is merged into `feat/local-infra-stack`):

| # | Slice | State |
|---|-------|-------|
| 1 | Local infrastructure stack | merged into `feat/local-infra-stack` (PR SanctumLabs/curtz#342, open) |
| 3 | App connectivity (config, health, shutdown, migrator) | merged into `feat/local-infra-stack` |
| 2 | **This spec** — image, `app` stack, CI | |
| 4 | OpenTelemetry | next |
| 5 | Outbox relay and the Kafka client | after 4 |

**In scope:** `Dockerfile`, `.dockerignore`, `hadolint.yaml`, a `healthcheck` subcommand in `app/cmd/main.go`, `deploy/app/compose.yml`
and the root `docker-compose.yml` include, `scripts/infra.sh` and `scripts/infra_test.sh`, `.make/docker.mk`, `.env.example` (only
if a variable is added), every file under `.github/workflows/`, `.github/dependabot.yml`, docs (`docs/LocalInfrastructure.md`,
`docs/Deployment.md`, `README.md`) and one ADR.

**Out of scope:** changing the secret names `deploy.yml` publishes (see §11), `fly.toml`, image signing, Kubernetes manifests,
OpenTelemetry, the Kafka client, the shared workflows in `SanctumLabs/ci-workflows`, and the broken `-tags legacy` build.

**Starting point.** After the pull into `feat/local-infra-stack` (tip `5d91b87`) the module did not compile: `server.go` used
`strconv` without importing it. That one-line fix is the first commit on this branch (`009f832`). The same pull made
`Server.Serve` honor `SERVER_HOST`, which this spec relies on (§5).

## 3. Decisions

| # | Decision | Why |
|---|----------|-----|
| D1 | The final stage is distroless `static` `:nonroot`, pinned by digest, with no shell. Confirmed by the user. | Smallest attack surface: no shell, no package manager, uid 65532, CA certificates and tzdata included. |
| D2 | The image carries `/app/curtz`, `/app/migrator` and `/app/migrations`. The compose `migrate` job stays on the `migrate/migrate` image. Confirmed by the user. | One artifact can serve and migrate (ADR-0014), and `make infra.postgres.up` still needs no app build. |
| D3 | CI is modernized as part of this slice. Confirmed by the user. | CI cannot build the app today. It can only be linted locally (`actionlint`); the first real run is on GitHub. |
| D4 | Health checking is a `curtz healthcheck` subcommand that calls `GET /health` (liveness). | A distroless image has no `curl` or `wget`. Liveness, not readiness, so a Postgres outage does not mark the container unhealthy. |
| D5 | The image sets `ENV ENVIRONMENT=production`. Compose overrides it with `development`. | A container started without `ENVIRONMENT` must refuse the development secrets instead of silently accepting them (slice 3 review finding). |
| D6 | The `app` stack is not part of `core-*` or `full-*`. | `make infra.core.up` and `make infra.full.up` must not trigger an image build. |
| D7 | The app container gets an explicit `environment:` list, not `env_file: .env`. | The `.env` file holds every stack's admin passwords. The app needs a handful of values. |
| D8 | Actions are pinned by commit SHA with a version comment, and Dependabot gets the `docker` and `github-actions` ecosystems. | Pins are only safe if something keeps them current. |
| D9 | `deploy.yml` passes secrets through `env:` and keeps its secret names. | Interpolating `${{ secrets.X }}` into a shell script is an injection risk. Which secrets exist in GitHub is the owner's decision. |
| D10 | Every commit made while implementing this spec ends with the `Co-Authored-By` trailer given in the session's attribution instructions, as the user asked on 2026-10-04. Example `git commit` snippets in plans omit the trailer text. | A project convention, stated here so it survives into the plan. |

## 4. The image

**Build arguments:** `VERSION` (default `unknown`), `GIT_COMMIT` (default `unknown`), `BUILD_TIME` (default `unknown`). They feed
the variables that `app/pkg/version.go` already declares.

**Stage `build`** (the Go toolchain image, an Alpine variant of Go 1.26, tag and digest resolved when implementing and recorded
in the implementation notes appended to this spec):

1. `COPY go.mod go.sum`, then `go mod download && go mod verify` with a BuildKit cache mount for `/go/pkg/mod`.
2. `COPY app ./app` (the allowlisted context, below).
3. `CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags "-s -w -X github.com/sanctumlabs/curtz/app/pkg.Version=$VERSION
   -X …GitCommit=$GIT_COMMIT -X …BuildTime=$BUILD_TIME"` for `./app/cmd` (to `/out/curtz`) and `./app/cmd/migrator` (to
   `/out/migrator`), with cache mounts for `/go/pkg/mod` and the build cache. There is no `-mod=mod`, so the build uses `go.sum` as is.

**Final stage** (distroless `static`, newest Debian release it publishes, `:nonroot`, pinned by digest):

- `COPY --from=build /out/curtz /out/migrator /app/` and `COPY app/internal/adapters/postgres/migrations /app/migrations`.
- `WORKDIR /app`, `USER 65532:65532`, `EXPOSE 8085`, `ENTRYPOINT ["/app/curtz"]`.
- `ENV ENVIRONMENT=production` (D5) and `ENV MIGRATIONS_PATH=/app/migrations`.
- `HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 CMD ["/app/curtz", "healthcheck"]`.
- OCI labels (`org.opencontainers.image.source`, `.revision`, `.version`, `.created`, `.title`, `.licenses`) from the build arguments.
- No `# syntax=` directive: it would pull a BuildKit frontend image, and the cache mounts work with Docker's built-in frontend.

**`.dockerignore`** becomes an allowlist: ignore everything, then allow `go.mod`, `go.sum` and `app/`, and drop `**/*_test.go` and
`app/test/`. A secret or a new top-level directory can then never enter the context by accident.

**`hadolint.yaml`:** add `gcr.io` to the trusted registries and drop the `DL3017`/`DL3018` ignores (the image no longer uses `apk`).

## 5. The `healthcheck` subcommand

`/app/curtz healthcheck` loads the server settings with `config.LoadServer` (so it honors `HTTP_PORT` and `SERVER_HOST`), calls
`GET http://<host>:<port>/health` with a 2 second timeout, prints one line, and exits 0 on HTTP 200 and 1 otherwise. A host of
`""`, `0.0.0.0` or `::` is dialed as `127.0.0.1`; any other configured host is dialed as given, because `Serve` now binds exactly
the configured host. It does not load `.env`, and it stays in `app/cmd/main.go` so `make run` (`go run app/cmd/main.go`) keeps
working. Tests, written first: success, a non-200 status, a refused connection, a timeout, and the host rewriting rule.

## 6. The `app` compose stack (`deploy/app/compose.yml`)

Included from the root `docker-compose.yml` with `env_file: .env`, like the other stacks.

- Services `app-ha` and `app-single`, profiles `app-ha` and `app-single`, both on the `curtz` network with the alias `curtz-app`.
  They build from the repository `Dockerfile` (`context: ../..`) and tag the result `curtz-app:local`.
- Environment (D7): `ENVIRONMENT=development`, `HTTP_PORT=8085`, `DATABASE_HOST=postgres`, `DATABASE_PORT=5432`,
  `DATABASE_NAME=${PG_DATABASE:-curtzdb}`, `DATABASE_USERNAME=${PG_APP_USER:-curtz-user}`,
  `DATABASE_PASSWORD=${PG_APP_PASSWORD:-curtz-pass}`, `REDIS_USERNAME=${REDIS_USERNAME:-curtz-svc}`,
  `REDIS_PASSWORD=${REDIS_PASSWORD:-curtz-svc}`, `AUTH_SECRET=${AUTH_SECRET:-curtz-secret}`, and `REDIS_ADDRESS` as
  `redis-1:7001` in single mode and `redis-1:7001` to `redis-6:7006` in HA mode. Every `${VAR:-default}` must match
  `.env.example`, which `scripts/infra_env_check.sh` enforces.
- Dependencies: none declared. A `depends_on` into a profile that is not enabled makes `docker compose config` fail (verified on
  Compose 5.5.1), so ordering is the script's job: `infra.sh up app` brings up and waits for Postgres (with its `migrate` job) and
  Redis first, and `restart: unless-stopped` covers a manual start before Postgres is ready. Redis is optional for the app (ADR-0015).
- Hardening: `read_only: true`, `tmpfs: [/tmp]`, `cap_drop: [ALL]`, `security_opt: [no-new-privileges:true]`,
  `restart: unless-stopped`, `stop_grace_period: 20s` (longer than the 15 second `SHUTDOWN_TIMEOUT`, so a stop never kills a drain).
- Port `127.0.0.1:8085:8085`, so stop a host-run API first. `SERVER_HOST` is not set: inside a container it must stay `0.0.0.0`.
- `scripts/infra.sh` learns the `app` stack (the mode applies): `up app <mode>` first brings up `postgres` and `redis` for that mode
  through the script itself, then starts the app profile with `--build`; `down app` stops only the app. `scripts/infra_test.sh`
  gets stub-based cases for it, and `make infra.config` must render both profiles.
- Make: `infra.app.up` and `infra.app.down` with `MODE=ha|single`. The global `infra.config`, `infra.ps`, `infra.logs` and `infra.clean`
  cover the app once `app-ha` and `app-single` are in `INFRA_PROFILES`.

## 7. Make, hadolint and Trivy

- `lint.docker` mounts `hadolint.yaml` with an absolute path (today it passes a relative name, which Docker treats as a named volume,
  so the config is never read) and pins the hadolint image by digest (the image already on the machine).
- `scan.docker` builds the image if missing, saves it with `docker save` to a temporary tar and scans that file with a pinned
  Trivy image (`--severity HIGH,CRITICAL --ignore-unfixed --exit-code 1`, a named volume for the vulnerability database). It does
  not mount the Docker socket into the scanner. The empty `scan.docker.image` target is removed.
- `build.docker` passes `VERSION` (`git describe --tags --always --dirty`), `GIT_COMMIT` and `BUILD_TIME`.
- A new `lint.workflows` target runs `actionlint` from a pinned image.

## 8. CI

| File | Change |
|------|--------|
| `docker.yml` | Rewritten. It checks out `workflow_run.head_sha`. It runs hadolint, builds once and scans with Trivy before anything is pushed, then builds `linux/amd64` and `linux/arm64` and pushes to Docker Hub and GHCR in one step with SBOM and provenance. It has least-privilege `permissions` and a GitHub Actions layer cache. The existing triggers and registries are unchanged. |
| `build_app.yml` | `go-version-file: go.mod` instead of `1.18.x`, `actions/upload-artifact@v4` (v2 and v3 no longer run), and the checkout fixed as above. |
| `tests.yml`, `lint.yml` | The Go version input goes from `1.25` to `1.26` to match `go.mod`. |
| `codeql.yml` | Current CodeQL action versions and checkout. |
| `release.yml`, `sentry.yml`, `dangerci.yml`, `slack.yml`, `bitbucket_sync.yml`, `gitlab_sync.yml` | Current action versions, Go from `go.mod` where it is set (`sentry.yml`), Node 22 instead of the end-of-life 16 (`dangerci.yml`). |
| `deploy.yml` | Current checkout, `setup-flyctl` pinned instead of `@master`, and secrets passed through `env:` (D9). Secret names unchanged. |
| every workflow | A minimal `permissions:` block, and actions pinned by SHA with a version comment (D8). |
| `dependabot.yml` | Adds the `docker` and `github-actions` ecosystems next to `gomod`. |

SHAs are resolved with `gh api` when implementing. If a pinned action's behavior changes (inputs renamed between major versions),
the workflow is adjusted to the current documented inputs.

## 9. Verification plan

Written test-first where there is code (the `healthcheck` subcommand, `infra.sh` cases); everything else is verified by running it.

- `make lint.docker`: no findings. `actionlint` on every workflow: no findings.
- `docker build` cold and cached; image size recorded; the build context shows no `.env` or `.git` (a debug build that lists `/src`).
- The container runs as uid 65532; `docker run --entrypoint sh …` fails (no shell); `--read-only --cap-drop ALL --security-opt
  no-new-privileges` still serves `/health/ready`; the Docker health status becomes `healthy`.
- `docker run --entrypoint /app/migrator …` against the stack: first run applies or "no change", second run "no change".
- `docker run -e ENVIRONMENT=production …` without secrets exits 1 naming the variables (the image default, D5).
- `make scan.docker`: no fixable HIGH or CRITICAL finding (any exception is listed with its reason in the implementation notes).
- `make infra.app.up MODE=single` and `MODE=ha`: `ready`, `healthy`, `/health/ready` ok, a registration returns 201; Redis stopped gives
  `degraded`; SIGTERM through `docker compose stop app-…` exits 0 inside the grace period; `make infra.app.down` removes only the app.
- `bash scripts/infra_test.sh`, `make infra.config` and `go test ./...` stay green.

## 10. Documentation

`docs/LocalInfrastructure.md` gains the `app` stack (commands, the port conflict with a host-run API, why there is no shell and how to
debug with the `:debug` distroless tag). `docs/Deployment.md` gains the image contract: the contents, the `ENVIRONMENT` default, the
healthcheck, the migrator entrypoint and the variables production must set. The README Docker paragraph points at both. ADR-0016
records D1, D2 and D5: the runtime image is distroless static, carries the migrator, and defaults to production.

## 11. Risks and open items

- **Downloads.** Building and verifying needs images that are not local: the Go toolchain image, the distroless image, Trivy (plus its
  vulnerability database) and `actionlint`. Each is pulled only with the user's go-ahead before the first build, as in slice 1.
- **CI is unproven until it runs on GitHub.** `actionlint` catches syntax and input errors, not runtime behavior. The multi-arch build
  adds emulated arm64 time to the workflow.
- **`deploy.yml` still publishes Mongo-era Fly secrets** (`DATABASE`, `CACHE_*`, `AUTH_EXPIRE`) that the app no longer reads. Production
  needs `DATABASE_*` (or `DATABASE_URL`), `REDIS_ADDRESS`, `REDIS_PASSWORD`, `AUTH_SECRET` and the rest of the contract in `.env.example`.
  Those secret names are the owner's call, and `fly.toml` (`ENV`, `PORT`) has the same mismatch. With D5 the app now refuses to start on
  Fly until they are set, which is the intended fail-fast.
- **No shell in the image.** Debugging uses the `:debug` distroless tag, documented in §10.
- **Digest pins go stale** without Dependabot; D8 adds it, but it only opens pull requests.
