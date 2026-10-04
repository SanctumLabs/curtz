# Deployment

The application ships as one container image that holds the API and the database migrator.

## Build

```bash
make build.docker DOCKER_IMAGE_TAG=curtz-service    # stamps the image with its version, commit and build time
# or
docker build -t <IMAGE_NAME>:<IMAGE_TAG> --build-arg VERSION=<VERSION> --build-arg GIT_COMMIT=<SHA> --build-arg BUILD_TIME=<RFC3339> .
```

The image is built in two stages. The final stage is `gcr.io/distroless/static:nonroot`, pinned by digest: no shell, no package manager, uid 65532. It contains:

- `/app/curtz`, the API (the default entrypoint, port 8085);
- `/app/migrator`, which applies the database migrations (ADR-0014);
- `/app/migrations`, the SQL files the migrator reads (`MIGRATIONS_PATH` is preset to it).

`make lint.docker` runs hadolint with the rules in `hadolint.yaml`, and `make scan.docker` scans the image for fixable HIGH and CRITICAL vulnerabilities. CI does both before it publishes (`.github/workflows/docker.yml`).

Without Docker: `make build` writes the binary to `bin/curtz`.

## Configuration

The image sets `ENVIRONMENT=production`. In that mode the API and the migrator refuse the development secrets and exit 1 naming the variables that are missing, so a deployment must set at least:

| Variable | Meaning |
|---|---|
| `AUTH_SECRET` | the JWT signing secret |
| `DATABASE_HOST`, `DATABASE_NAME`, `DATABASE_USERNAME`, `DATABASE_PASSWORD` | Postgres (the primary, port 5432); or `DATABASE_URL` instead, which wins over those |
| `DATABASE_SSL_MODE` | defaults to `disable`; use `require` or stronger outside a private network (the API logs a warning otherwise) |
| `REDIS_ADDRESS`, `REDIS_PASSWORD` | a comma-separated `host:port` list (one address gives a plain client, several a cluster client) and the password; `REDIS_USERNAME` defaults to `curtz-svc` |

Every variable, its default and its unit is listed in the app-connectivity spec (`docs/superpowers/specs/2026-10-04-app-connectivity-design.md`, section 4); the local defaults are in `.env.example`.

## Migrations

The API never migrates at startup. Run the migrator from the same image, once per deploy, before starting the new version:

```bash
docker run --rm -e ENVIRONMENT=production -e DATABASE_HOST=... -e DATABASE_NAME=... -e DATABASE_USERNAME=... -e DATABASE_PASSWORD=... \
  --entrypoint /app/migrator <IMAGE_NAME>:<IMAGE_TAG>
```

It exits 0 when it applied the migrations or found nothing to do. If it reports a connection error because the database is still starting, run it again.

## Running

```bash
docker run -p 8085:8085 -e AUTH_SECRET=... -e DATABASE_HOST=... -e DATABASE_PASSWORD=... -e REDIS_ADDRESS=... -e REDIS_PASSWORD=... \
  --name curtz-api <IMAGE_NAME>:<IMAGE_TAG>
```

- `GET /health` is liveness (always 200 while the process runs); `GET /health/ready` is readiness (503 when Postgres is down or the process is draining).
- The image's `HEALTHCHECK` runs `/app/curtz healthcheck`, which calls `/health`.
- On SIGTERM the API stops accepting traffic and finishes in-flight requests for up to `SHUTDOWN_TIMEOUT` seconds (default 15). Give the orchestrator a stop grace period longer than that (the local compose stack uses 20 seconds).
- For hardening, also run it with a read-only root filesystem, `--cap-drop ALL` and `--security-opt no-new-privileges`, as `deploy/app/compose.yml` does.

## Continuous integration

GitHub Actions runs the full pipeline. GitLab CI (`.gitlab-ci.yml`) and Bitbucket Pipelines (`bitbucket-pipelines.yml`) run the same verification on the mirrors that `.github/workflows/gitlab_sync.yml` and `bitbucket_sync.yml` keep in step:

| Check | Command | GitHub | GitLab | Bitbucket |
|---|---|---|---|---|
| Lint | `golangci-lint run` (v2.13.2, the version `make lint` pins) | yes | yes | yes |
| Unit tests with coverage | `make test.coverage` | yes | yes | yes |
| Integration tests | `make test.integration` (testcontainers) | yes | yes | yes |
| End-to-end tests | `make test.e2e` (testcontainers) | yes | yes | yes |
| Binary build | `make build` | Linux, macOS, Windows | Linux | Linux |
| Dockerfile lint | `hadolint --config hadolint.yaml Dockerfile` | yes | yes | yes |
| Image build and checks | `docker build`, then `scripts/image_test.sh image <tag>` | yes | yes | yes |
| Image scan | `trivy image --severity HIGH,CRITICAL --ignore-unfixed --exit-code 1` | yes | yes | yes |

Publishing the image, releases, the Fly.io deploy, Sentry releases, Slack notifications, CodeQL and Danger run on GitHub only.

The jobs that need Docker (the integration and end-to-end tests and the image build) use Docker-in-Docker on GitLab (`docker:29.8-dind`) and the `docker` service on Bitbucket, with `TESTCONTAINERS_RYUK_DISABLED=true`. The image is scanned from a saved tar, so the scanner needs no Docker daemon. `make lint.workflows` runs `scripts/workflows_check.sh`, `actionlint` and `scripts/mirror_ci_check.sh`, which fails when a mirror drifts from the GitHub workflows (a Go version that differs from `go.mod`, a missing command, an unpinned image). It can only check the files, not run them: the first run of a mirror pipeline happens when GitHub pushes to it.

## Local infrastructure

The supporting services (Postgres, Redis, Kafka, ELK, Prometheus, Grafana) run locally in Docker in either HA or single-node mode, and the API can run beside them as a container (`make infra.app.up`). See [Local infrastructure](./LocalInfrastructure.md).
