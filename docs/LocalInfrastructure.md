# Local infrastructure

Curtz depends on Postgres, Redis and Kafka, and is observed with ELK (logs), Prometheus (metrics), Tempo (traces) and
Grafana. This repository runs all of them locally with Docker Compose. Each stack runs in one of two modes:

- **HA**: the clustered topology used in production (3 Kafka brokers, a 6-node Redis Cluster, Postgres with Patroni and
  2 replicas, 3 Elasticsearch nodes). Use it to exercise failover.
- **single**: one node per stack. Use it for everyday work; it needs a fraction of the memory.

Nothing starts without a profile: a bare `docker compose up` does nothing. Use the `make infra.*` commands below.

## Prerequisites

- Docker with Compose v2.20 or newer (`docker compose version`) and `make`.
- Docker Desktop memory (Settings, Resources): see the [memory budget](#memory-budget). HA stacks need much more than single mode.
- Linux hosts running ELK: `sudo sysctl -w vm.max_map_count=262144` (Docker Desktop already sets it).
- Verified on Docker Desktop for Mac (Apple silicon). Linux, rootless Docker, SELinux labels and non-default log drivers have not
  been exercised; Filebeat in particular reads `/var/lib/docker/containers` and assumes the default `json-file` log driver.
- Every `make infra.*` target creates `.env` from `.env.example` when it is missing. Raw `docker compose` commands need
  `.env` to exist, so run `make create.envfile` first.

## Quick start

```bash
make infra.core.up MODE=single      # Postgres + Redis + Kafka, one node each (lightest)
make infra.core.up                  # the same, HA (3 Postgres, 6 Redis, 3 Kafka)
make infra.observability.up         # Prometheus, Grafana, Tempo, Alertmanager, OpenTelemetry Collector
make infra.elk.up MODE=single       # Elasticsearch, Logstash, Kibana, Filebeat
make infra.full.up MODE=single      # every stack except legacy
make infra.ps                       # what is running and whether it is healthy
make infra.core.down                # stop (data is kept)
make infra.clean                    # remove every container AND volume except the legacy MongoDB and Redis (asks first)
```

`up` waits until every service is healthy and every one-shot job (topic creation, cluster formation, migrations,
Elasticsearch setup) has finished, or fails after five minutes and prints what is still pending.

## Modes and profiles

| Stack | HA profile | Single profile | Make target |
|---|---|---|---|
| Kafka + Kafka UI | `kafka-ha` | `kafka-single` | `infra.kafka.up` |
| Redis | `redis-ha` | `redis-single` | `infra.redis.up` |
| Postgres + migrations | `postgres-ha` | `postgres-single` | `infra.postgres.up` |
| ELK | `elk-ha` | `elk-single` | `infra.elk.up` |
| Postgres + Redis + Kafka | `core-ha` | `core-single` | `infra.core.up` |
| Everything except legacy | `full-ha` | `full-single` | `infra.full.up` |
| Prometheus, Grafana, Tempo, Alertmanager, Collector | `observability` | `observability` | `infra.observability.up` |
| MongoDB + standalone Redis (deprecated) | `legacy` | `legacy` | `infra.legacy.up` |

- Stacks are independent: Kafka in HA with Postgres single is fine.
- **Never run both modes of one stack.** They use the same ports and names. `make infra.<stack>.up MODE=...` removes the
  other mode's containers first (volumes are kept, so switching back restores the data). With raw `docker compose`,
  stop one mode before starting the other.
- Observability has no HA mode: Prometheus, Grafana, Tempo and the Collector are single instances in every mode.

## Make commands

| Command | What it does |
|---|---|
| `make infra.<stack>.up [MODE=ha\|single]` / `.down` | start / stop a stack: `kafka redis postgres elk observability legacy core full` |
| `make infra.wait STACK=kafka MODE=single` | wait until a stack is healthy |
| `make infra.ps` / `infra.stats` | status and health / live memory and CPU |
| `make infra.logs SERVICE=kafka-1` | follow logs (omit `SERVICE` for everything) |
| `make infra.config` | check that every profile renders and the env defaults agree |
| `make infra.clean` | remove all containers and volumes except the legacy MongoDB and Redis data |
| `make infra.clean.legacy` | remove the legacy MongoDB and Redis containers and their volumes (their data is lost) |
| `make infra.hosts` | print the `/etc/hosts` line needed for Redis HA from an app on your host |
| `make infra.migrate` | re-run the database migrations |
| `make infra.psql` | `psql` as the application user on the primary |
| `make infra.redis.cli` | `redis-cli` as the application user, cluster mode |
| `make infra.kafka.topics` | describe the Kafka topics |
| `make infra.patroni.list` | Patroni members and roles (Postgres HA) |
| `make infra.es.health` | Elasticsearch cluster health |

The `MODE` variable also selects which node the helper commands talk to (for example `make infra.psql MODE=single`).

## Connecting the application

The same names work in both modes; the lists just get shorter. Use the "host" column when the app runs on your machine
(for example `make run.dev`) and the "network" column when it runs in a container on the `curtz` network.

| Concern | App on the host | App in the compose network |
|---|---|---|
| Kafka bootstrap | HA `localhost:19092,localhost:29092,localhost:39092`; single `localhost:19092` | HA `kafka-1:9092,kafka-2:9092,kafka-3:9092`; single `kafka-1:9092` |
| Redis seed nodes | HA `localhost:7001`..`localhost:7006`; single `localhost:7001` | HA `redis-1:7001`..`redis-6:7006`; single `redis-1:7001` |
| Postgres (write) | `localhost:5432` | `postgres:5432` |
| Postgres (read) | `localhost:5433` | HA `postgres:5433`; single `postgres:5432` |
| OTLP (traces, metrics) | `http://localhost:4317` (gRPC) or `:4318` (HTTP) | `http://otel-collector:4317` |

Notes:

- **Redis HA announces hostnames** (`redis-1`..`redis-6`) so cluster clients can follow redirects. An app on your host must
  resolve them: run `make infra.hosts` and add the printed line to `/etc/hosts` once. A containerised app needs nothing.
  Single mode is a cluster of one, so it behaves like the HA cluster for multi-key rules (`CROSSSLOT`).
- Redis serves only logical database `0` (cluster mode), so `REDIS_DATABASE` must be `0`.
- Postgres reads on `5433` reach a replica in HA (and the same node in single mode); writes on a read port fail with
  "read-only transaction".
- Logs need no configuration: write JSON to stdout and Filebeat ships it (see ELK).

## Running the app against the stack

The API defaults match the stack, so on your host no `.env` edits are needed.

```bash
make infra.core.up MODE=single      # or HA: also run `make infra.hosts` once, and use the six-address REDIS_ADDRESS in .env.example
go run ./app/cmd/migrator           # optional: infra.core.up already migrated through the compose job; a second run says "no change"
make run                            # the API on :8085
curl -s localhost:8085/health/ready
```

`go run ./app/cmd/migrator` is the same migration code the tests use (`postgres.Migrate`). Run it from the repository root, or set `MIGRATIONS_PATH`. `make run.with.migrations` runs it and then the API. The API never migrates at startup (ADR-0014).

| Endpoint | Meaning |
|---|---|
| `GET /health` | liveness: 200 while the process runs |
| `GET /health/ready` | readiness: 200 `ok` (everything up) or `degraded` (Redis down); 503 `unavailable` (Postgres down) or `draining` (shutting down) |

Readiness never includes error text; look in the API's log for the reason (ADR-0015). Postgres is required, Redis is optional, and the API does not connect to Kafka.

On SIGTERM or Ctrl-C the API turns readiness to 503, finishes in-flight requests for up to `SHUTDOWN_TIMEOUT` seconds (default 15), closes Redis and Postgres and exits 0. A second Ctrl-C during the drain ends it at once.

The variables are listed, with their defaults, in `.env.example`. A value that does not parse stops startup with a message naming the variable, and any `ENVIRONMENT` other than `development` or `test` refuses the development secrets (`AUTH_SECRET`, `DATABASE_PASSWORD`, `REDIS_PASSWORD`).

Debugging:

- `redis is down, continuing without it` at startup: check `REDIS_ADDRESS` (HA needs all six seed nodes and the `/etc/hosts` line from `make infra.hosts`) and `REDIS_USERNAME`/`REDIS_PASSWORD`.
- Postgres connection errors: the write port is `5432` in both modes; `5433` is the HA read port and rejects writes.
- An old `.env` copied from before this change carries Mongo-era values (`DATABASE_PORT=27017`, `REDIS_ADDRESS=localhost`,
  `AUTH_SECRET=<AUTH_SECRET>`, ...). They override the new defaults, so the API or the migrator dials the wrong port or stops
  with an invalid-configuration error. Refresh it with `cp .env.example .env` (re-apply any values of your own first).
- `make infra.psql`, `make infra.redis.cli` and `make infra.patroni.list` show the other side of each connection.

## The stacks

### Kafka

- HA: `kafka-1..3` in KRaft mode, each both broker and controller; replication factor 3, `min.insync.replicas=2`.
  Single: `kafka-single` (alias `kafka-1`), replication factor 1.
- Topics are created at start-up and are idempotent: `url.events`, `identity.events`, `url.access`, `url.invalidations`,
  `security.scan.requests`, `webhook.deliveries`. Auto-creation is off.
- Kafka UI: <http://localhost:8080>. Describe topics with `make infra.kafka.topics`.
- No authentication or TLS locally.

### Redis

- HA: `redis-1..6` (3 masters, 3 replicas). Single: `redis-single`, one node owning all slots. Both run in cluster mode.
- Eviction `allkeys-lru`, 128 MB per node, AOF persistence. Application user `curtz-svc` cannot run dangerous commands
  (`FLUSHALL`, `KEYS`, `CONFIG`...).
- `make infra.redis.cli` opens a cluster-aware shell.

### Postgres

- HA: `patroni-1..3` managed by Patroni, a 3-member etcd quorum and HAProxy. `localhost:5432` always reaches the current
  primary, `localhost:5433` the replicas. Replication is asynchronous; set `synchronous_mode: true` in
  `deploy/postgres/patroni.yml` for zero data loss at the cost of write latency.
- Single: one plain Postgres 18 node; both ports reach it.
- The `migrate` job applies `app/internal/adapters/postgres/migrations` on every start (a no-op when up to date) and
  records them in `schema_migrations`.
- HAProxy stats: <http://localhost:8404>. Patroni REST: `localhost:8008`..`8010` (`/cluster`, `/metrics`).

### ELK

- Path: your app logs JSON to stdout, Filebeat reads every container's logs of this compose project, Logstash parses
  them, Elasticsearch stores them in the data stream `logs-curtz-default` (rolled over daily, deleted after 7 days),
  Kibana shows them.
- JSON log lines are mapped to `message`, `log.level`, `service.name`, `trace.id`, `span.id`; any other line keeps its text
  and takes the compose service as `service.name`.
- HA: 3 Elasticsearch nodes (transport TLS between them), 2 Logstash, 2 Kibana behind nginx. Single: one of each.
- Kibana: <http://localhost:5601>, log in as `elastic`. First time: Stack Management, Data Views, create `logs-curtz-*`
  with timestamp `@timestamp`.
- Filebeat runs as root and mounts the Docker socket read-only to read container logs. This is the one least-privilege
  exception, acceptable for local development only.

### Observability

- The app sends OTLP to the OpenTelemetry Collector (`localhost:4317`). Traces go to Tempo, metrics are re-exposed to
  Prometheus. Logs do not pass through the Collector.
- Grafana: <http://localhost:3000>. Datasources Prometheus, Tempo and Elasticsearch are provisioned; Tempo links spans to
  logs by `trace.id`. Dashboards "Stack overview" and "Curtz service" are in the folder "Curtz". The service dashboard
  shows "No data" until the app emits OpenTelemetry metrics.
- Prometheus: <http://localhost:9090>; Alertmanager: <http://localhost:9093>; Tempo: <http://localhost:3200>.
- Prometheus finds exporters by DNS, so only running components appear as targets. To get notifications, replace the
  `null` receiver in `deploy/observability/alertmanager.yml` with a Slack, email or webhook receiver and reload.
- More dashboards: import community dashboards by ID in Grafana (Dashboards, New, Import): Kafka 7589, Redis 763,
  Postgres 9628, Elasticsearch 14191, Go processes 6671.

### Legacy

`make infra.legacy.up` starts MongoDB 4.4 (end-of-life) and the old standalone Redis for the code behind the `legacy`
build tag. They are unchanged apart from binding to `127.0.0.1`.

## Credentials and ports

All credentials are development defaults from `.env.example`. You can override them in `.env`, with two rules:

- Values may only contain letters, digits, `.`, `_` and `-`. They end up in shell scripts, SQL, JSON bodies and connection
  URLs, where anything else (`$`, `#`, `@`, `/`, `:`, `%`, quotes, spaces) silently changes the meaning. `make infra.<stack>.up`
  and `make infra.config` refuse to continue when a value breaks this.
- Passwords are applied when a volume is first created (Postgres roles, the Elasticsearch `elastic` user, Grafana's admin).
  Changing one later does not change the stored password, and the stack then fails to authenticate. Run `make infra.clean`
  (all data is lost) or change the password inside the service as well.

| Service | Address | Login |
|---|---|---|
| Kafka (host) | `localhost:19092` (HA also `29092`, `39092`) | none |
| Kafka UI | <http://localhost:8080> | none |
| Redis | `localhost:7001`..`7006` | user `curtz-svc`, password `curtz-svc`; admin (default user) `curtz-redis-admin` |
| Postgres | `localhost:5432` write, `5433` read | `curtz-user` / `curtz-pass`, database `curtzdb`; superuser `postgres` / `curtz-postgres-admin` |
| Elasticsearch | `localhost:9200` | `elastic` / `curtz-elastic-dev` |
| Kibana | <http://localhost:5601> | `elastic` / `curtz-elastic-dev` |
| Grafana | <http://localhost:3000> | `admin` / `curtz-grafana-dev` |
| Prometheus, Alertmanager, Tempo | `9090`, `9093`, `3200` | none |
| OTLP | `4317` gRPC, `4318` HTTP | none |
| Legacy MongoDB / Redis | `27017` / `6379` | `curtzUser` / `curtzPassword` / none |

Every published port is bound to `127.0.0.1`. Do not reuse these values anywhere else.

## Memory budget

Docker Desktop has a fixed memory allocation (Settings, Resources). If a container dies with exit code 137 it was killed
for lack of memory. These figures were measured with `make infra.stats` on idle stacks (Docker Desktop on Apple silicon,
7.75 GiB allocated); expect more once Elasticsearch holds data. `full-ha` is the sum of the measured parts.

| Profile | Single | HA |
|---|---|---|
| `core` (Postgres, Redis, Kafka, Kafka UI) | 0.75 GiB | 2.0 GiB |
| `elk` | 3.0 GiB | 6.2 GiB |
| `observability` | 0.5 GiB | 0.5 GiB |
| `full` | **4.2 GiB** | **8.7 GiB**: does not fit 7.75 GiB, raise Docker to at least 12 GiB |

`full-single` therefore runs comfortably in the default allocation, and `elk-ha` alone needs most of it. Per-container
numbers: `make infra.stats`.

Heaps are small on purpose (512 MB for Kafka, Elasticsearch and Logstash); tune them with `KAFKA_HEAP`, `ES_HEAP` and
`LS_HEAP` in `.env`. On a small allocation run one stack at a time.

## Failure drills (HA)

| Stack | Do | Expect |
|---|---|---|
| Kafka | `docker compose --profile '*' stop kafka-2` | produce and consume keep working; `start` it and under-replicated partitions drain |
| Redis | stop a master (see `cluster nodes` in `make infra.redis.cli`) | a replica is promoted within ~10 seconds; writes continue |
| Postgres | stop the leader shown by `make infra.patroni.list` | a replica becomes leader (a few seconds after a graceful stop, about 30 seconds after a crash, measured); `localhost:5432` follows it once HAProxy's health checks see the new primary (about 6 more seconds) |
| Elasticsearch | stop `es-2` | cluster stays green or yellow; logs keep arriving |
| Logstash / Kibana | stop `logstash-1` / `kibana-1` | Filebeat uses `logstash-2`; Kibana stays reachable through nginx |

## Troubleshooting

First look: `make infra.ps` (health), `make infra.logs SERVICE=<name>`, `make infra.stats` (memory). With raw
`docker compose`, add `--profile '*'` to address a service by name.

| Symptom | Likely cause and fix |
|---|---|
| `stat .../.env: no such file` | raw `docker compose` without `.env`: `make create.envfile` |
| container `Exited (137)` | out of memory: raise Docker's memory, lower `*_HEAP`, or run fewer stacks |
| `port is already allocated` | the other mode of that stack is running, or another program uses the port (`lsof -i :5432`); `make infra.<stack>.up` clears the other mode |
| `no such service: kafka-1` | add `--profile '*'` to the raw `docker compose` command |
| `make infra.<stack>.up` times out | the message lists the pending services; read their logs |
| Kafka client from the host cannot connect | use the `localhost:19092,...` addresses; inside the network use `kafka-1:9092` |
| Kafka produce fails with `NOT_ENOUGH_REPLICAS` | HA needs 2 of 3 brokers for `acks=all`; start the stopped broker |
| Redis `lookup redis-2: no such host` | app on the host without the hosts entry: `make infra.hosts` |
| `make infra.<stack>.up` stops with "may only contain letters, digits..." | a value in `.env` has a character outside `A-Za-z0-9._-`; fix it there |
| Authentication fails after changing a password in `.env` | passwords are fixed when a volume is first created: `make infra.clean` (data is lost) or change it in the service too |
| Redis `CLUSTERDOWN` | wait ~10 seconds after a failure; check `cluster info`; if volumes were partly wiped, `make infra.clean` and start again |
| Redis `CROSSSLOT` | a multi-key command spans slots; use hash tags (`{user42}:a`). This is the same in production |
| Postgres connection refused right after start (HA) | no primary elected yet: `make infra.patroni.list`, wait for a Leader |
| Postgres `cannot execute ... in a read-only transaction` | writing to port 5433; use 5432 |
| Patroni member stuck, wrong timeline | `docker compose --profile '*' exec patroni-1 /opt/patroni/bin/patronictl -c /etc/patroni/patroni.yml reinit curtz patroni-2` |
| Migration says dirty | fix the SQL, then `docker run --rm --network curtz -v "$PWD/app/internal/adapters/postgres/migrations:/m" migrate/migrate -path=/m -database "postgres://curtz-user:curtz-pass@postgres:5432/curtzdb?sslmode=disable&x-migrations-table=schema_migrations" force <version>` |
| Elasticsearch `max virtual memory areas vm.max_map_count` (Linux) | `sudo sysctl -w vm.max_map_count=262144` |
| Kibana "server is not ready" | it needs 1 to 2 minutes after Elasticsearch is healthy |
| Kibana shows no logs | create the `logs-curtz-*` data view; check `make infra.logs SERVICE=filebeat-single` (or `filebeat-ha`); only containers of the `curtz` compose project are collected |
| Grafana Elasticsearch datasource error | ELK is not running |
| An exporter target is missing in Prometheus | its stack is not running; targets are discovered by DNS |

## Local versus production

| Topic | Locally | Production |
|---|---|---|
| Kafka auth/TLS | none | SASL and TLS |
| Elasticsearch HTTP | plain HTTP, auth on | HTTPS |
| Postgres replication | asynchronous | consider synchronous |
| Redis hostnames | via `/etc/hosts` | real DNS |
| Prometheus, Grafana, Tempo | single instance | replicated (Thanos or Mimir for Prometheus) |
| Filebeat | root, Docker socket | node agent with least privilege |
| Credentials | committed defaults | secrets manager |

## Adding a service

1. Put it in the stack's file under `deploy/<stack>/compose.yml` (or a new `deploy/<name>/compose.yml` plus an
   `include:` entry in `docker-compose.yml`).
2. Give it profiles: its own `<stack>-ha` and/or `<stack>-single`, `core-*` if the app needs it, `full-*`. For a mode-specific
   variant, give the single service the alias of the HA node 1 so addresses do not change.
3. Bind published ports to `127.0.0.1`, pin the image tag, and give every credential `${VAR:-default}` with the same
   default in `.env.example`.
   Name one-shot jobs `*-init-*`, `*-setup-*` or `migrate` (with `restart: "no"` or `on-failure`): `make infra.<stack>.up`
   waits for them to exit 0, while every other service must be running and healthy.
4. Run `make infra.config`.
