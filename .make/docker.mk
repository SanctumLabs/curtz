############################################################################################################################################################################################################
## Docker scripts
############################################################################################################################################################################################################
# make arguments that are defaulted
DOCKER_FILE ?= Dockerfile
DOCKER_IMAGE_TAG ?= curtz-service

.PHONY: create.dockerenvfile
create.dockerenvfile: ## Create a docker environment file
	if [ ! -f .env.docker ]; then cp .env.example .env.docker; fi

scan.docker.image:
	@echo "Scanning Docker Image: $(IMAGE)"
	./bin/trivy $(IMAGE)

# See local hadolint install instructions: https://github.com/hadolint/hadolint
.PHONY: lint.docker
lint.docker: ## lints the Dockerfile
	@echo "Running lint checks on Dockerfile"
	docker run --rm -i -v hadolint.yaml:/.config/hadolint.yaml hadolint/hadolint < $(DOCKER_FILE)
	@echo "Done linting Dockerfile"

# Reference: https://trivy.dev/latest/getting-started/
.PHONY: scan.docker
scan.docker: ## scans a docker image for vulnerabilities, but first it will build the image
	@if ! docker image inspect $(DOCKER_IMAGE_TAG) >/dev/null 2>&1; then \
		echo ">>> Building Docker image '$(DOCKER_IMAGE_TAG)' as it does not exist locally <<<<"; \
		docker build -f $(DOCKER_FILE) . -t $(DOCKER_IMAGE_TAG); \
		echo ">>> Done building Docker image $(DOCKER_IMAGE_TAG), scanning for vulnerabilities <<<<"; \
		docker run -v /var/run/docker.sock:/var/run/docker.sock -v ~/Library/Caches:/root/.cache/ aquasec/trivy image $(DOCKER_IMAGE_TAG); \
	else \
		echo ">>> Scanning Docker $(DOCKER_IMAGE_TAG) image for vulnerabilities <<<<"; \
		docker run -v /var/run/docker.sock:/var/run/docker.sock -v ~/Library/Caches:/root/.cache/ aquasec/trivy image $(DOCKER_IMAGE_TAG); \
		echo "\n >>> Done scanning docker image $(DOCKER_IMAGE_TAG) for vulnerabilities"; \
	fi

.PHONY: build.docker
build.docker: ## Build Docker image
	@echo "Building Docker image"
	docker build -f $(DOCKER_FILE) . -t $(DOCKER_IMAGE_TAG)
	@echo "Done building Docker image"

.PHONY: push.docker
push.docker: ## Pushes Docker image if it exists, otherwise it builds it, usage: 'make push.docker DOCKER_IMAGE_TAG=notification-svc', the DOCKER_IMAGE_TAG is optional
	@if ! docker image inspect $(DOCKER_IMAGE_TAG) >/dev/null 2>&1; then \
		echo "Docker Image $(DOCKER_IMAGE_TAG) does not exist, building..."; \
		docker build -f $(DOCKER_FILE) . -t $(DOCKER_IMAGE_TAG); \
		echo "Done building docker Image $(DOCKER_IMAGE_TAG), pushing image"; \
		echo "Pushing Docker image $(DOCKER_IMAGE_TAG)"; \
		docker push $(DOCKER_IMAGE_TAG); \
		echo "Done pushing docker image $(DOCKER_IMAGE_TAG)"; \
	else \
		echo "Pushing Docker image $(DOCKER_IMAGE_TAG)"; \
		docker push $(DOCKER_IMAGE_TAG); \
		echo "Done pushing docker image $(DOCKER_IMAGE_TAG)"; \
	fi

.PHONY: run.docker
run.docker: create.dockerEnvFile # Run the Docker container
	@if ! docker image inspect $(DOCKER_IMAGE_TAG) >/dev/null 2>&1; then \
		echo "Building Docker Image $(DOCKER_IMAGE_TAG)"; \
		docker build -f $(DOCKER_FILE) . -t $(DOCKER_IMAGE_TAG); \
		echo "Done building docker Image $(DOCKER_IMAGE_TAG), running container"; \
		docker run --env-file .env.docker -t $(DOCKER_IMAGE_TAG); \
	else \
		echo "Running docker image $(DOCKER_IMAGE_TAG)"; \
		docker run --env-file .env.docker -t $(DOCKER_IMAGE_TAG); \
	fi

# Ref: https://github.com/slimtoolkit/slim
.PHONY: slim.dockerxray
slim.dockerxray: ## Runs an xray command using slim on the docker image, usage: 'make slim.dockerxray DOCKER_IMAGE_TAG=notification-svc', the DOCKER_IMAGE_TAG is optional
	@if ! docker image inspect $(DOCKER_IMAGE_TAG) >/dev/null 2>&1; then \
		echo ">>> Building Docker image '$(DOCKER_IMAGE_TAG)' as it does not exist locally <<<<"; \
		docker build -f $(DOCKER_FILE) . -t $(DOCKER_IMAGE_TAG); \

		echo ">>> Done building Docker image $(DOCKER_IMAGE_TAG), running Xray <<<<"; \

		docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock dslim/slim xray $(DOCKER_IMAGE_TAG); \
	else \
		echo ">>> Running Xray Docker image '$(DOCKER_IMAGE_TAG)' <<<<"; \
		docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock dslim/slim xray $(DOCKER_IMAGE_TAG); \
	fi


############################################################################################################################################################################################################
## Local infrastructure (docker compose, see docs/LocalInfrastructure.md)
############################################################################################################################################################################################################
# MODE picks the topology for stacks that have one: ha (default) or single, e.g. make infra.kafka.up MODE=single
MODE ?= ha
INFRA := $(ROOT_DIR)/scripts/infra.sh
COMPOSE := docker compose --project-directory $(ROOT_DIR)
INFRA_PROFILES := kafka-ha kafka-single redis-ha redis-single postgres-ha postgres-single elk-ha elk-single observability legacy core-ha core-single full-ha full-single

# Service that answers for "node 1" of a stack in the selected mode, used by the debugging helpers
KAFKA_SERVICE := $(if $(filter single,$(MODE)),kafka-single,kafka-1)
REDIS_SERVICE := $(if $(filter single,$(MODE)),redis-single,redis-1)
POSTGRES_SERVICE := $(if $(filter single,$(MODE)),postgres-single,patroni-1)
ES_SERVICE := $(if $(filter single,$(MODE)),es-single,es-1)

.PHONY: infra.kafka.up infra.kafka.down
infra.kafka.up: create.envfile ## Start Kafka and Kafka UI (MODE=ha or single)
	@$(INFRA) up kafka $(MODE)
infra.kafka.down: ## Stop Kafka (data is kept)
	@$(INFRA) down kafka

.PHONY: infra.redis.up infra.redis.down
infra.redis.up: create.envfile ## Start the Redis cluster (MODE=ha or single)
	@$(INFRA) up redis $(MODE)
infra.redis.down: ## Stop Redis (data is kept)
	@$(INFRA) down redis

.PHONY: infra.postgres.up infra.postgres.down
infra.postgres.up: create.envfile ## Start Postgres and run migrations (MODE=ha or single)
	@$(INFRA) up postgres $(MODE)
infra.postgres.down: ## Stop Postgres (data is kept)
	@$(INFRA) down postgres

.PHONY: infra.elk.up infra.elk.down
infra.elk.up: create.envfile ## Start Elasticsearch Logstash Kibana and Filebeat (MODE=ha or single)
	@$(INFRA) up elk $(MODE)
infra.elk.down: ## Stop ELK (data is kept)
	@$(INFRA) down elk

.PHONY: infra.observability.up infra.observability.down
infra.observability.up: create.envfile ## Start Prometheus Grafana Tempo Alertmanager and the OTel Collector
	@$(INFRA) up observability
infra.observability.down: ## Stop the observability stack (data is kept)
	@$(INFRA) down observability

.PHONY: infra.legacy.up infra.legacy.down
infra.legacy.up: create.envfile ## Start the legacy MongoDB and standalone Redis
	@$(INFRA) up legacy
infra.legacy.down: ## Stop the legacy services (data is kept)
	@$(INFRA) down legacy

.PHONY: infra.core.up infra.core.down
infra.core.up: create.envfile ## Start everything the app needs - Postgres Redis Kafka (MODE=ha or single)
	@$(INFRA) up core $(MODE)
infra.core.down: ## Stop Postgres Redis and Kafka (data is kept)
	@$(INFRA) down core

.PHONY: infra.full.up infra.full.down
infra.full.up: create.envfile ## Start every stack except legacy (MODE=ha or single)
	@$(INFRA) up full $(MODE)
infra.full.down: ## Stop every stack except legacy (data is kept)
	@$(INFRA) down full

.PHONY: infra.config
infra.config: create.envfile ## Check that every compose profile renders and the env defaults agree
	@$(ROOT_DIR)/scripts/infra_env_check.sh
	@for p in $(INFRA_PROFILES); do \
		$(COMPOSE) --profile $$p config --quiet || exit 1; \
		echo "ok: $$p"; \
	done
	@$(COMPOSE) --profile '*' config --quiet && echo "ok: all profiles"

.PHONY: infra.ps
infra.ps: ## Show every infrastructure container and its health
	@$(COMPOSE) --profile '*' ps -a

.PHONY: infra.logs
infra.logs: ## Follow logs, usage - make infra.logs SERVICE=kafka-1 (omit SERVICE for everything)
	@$(COMPOSE) --profile '*' logs -f --tail=100 $(SERVICE)

.PHONY: infra.stats
infra.stats: ## Show memory and CPU of the running infrastructure containers
	@docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}\t{{.CPUPerc}}" \
		$$(docker ps --filter label=com.docker.compose.project=curtz -q)

.PHONY: infra.clean
infra.clean: confirm ## Remove every infrastructure container AND volume (all data is lost)
	@$(COMPOSE) --profile '*' down -v --remove-orphans

.PHONY: infra.hosts
infra.hosts: ## Print the hosts-file line needed to run the app on the host against Redis HA
	@echo "Redis HA announces the hostnames redis-1 to redis-6. An app running on your host must resolve them."
	@echo "Add this line to /etc/hosts (needs sudo, not required when the app runs in a container):"
	@echo "  127.0.0.1 redis-1 redis-2 redis-3 redis-4 redis-5 redis-6"
	@if grep -q 'redis-1' /etc/hosts; then echo "(an entry for redis-1 already exists)"; fi


.PHONY: infra.kafka.topics
infra.kafka.topics: ## Describe the Kafka topics (MODE=ha or single)
	@$(COMPOSE) --profile '*' exec -T $(KAFKA_SERVICE) /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --describe


.PHONY: infra.redis.cli
infra.redis.cli: ## Open redis-cli as the application user in cluster mode (MODE=ha or single)
	@$(COMPOSE) --profile '*' exec $(REDIS_SERVICE) sh -c 'redis-cli -p "$$NODE_PORT" --user "$$REDIS_USERNAME" --pass "$$REDIS_PASSWORD" --no-auth-warning -c'


.PHONY: infra.psql
infra.psql: ## Open psql as the application user on the primary (MODE=ha or single)
	@$(COMPOSE) --profile '*' exec $(POSTGRES_SERVICE) sh -c 'PGPASSWORD="$$PG_APP_PASSWORD" psql -h postgres -U "$$PG_APP_USER" "$$PG_DATABASE"'

.PHONY: infra.migrate
infra.migrate: create.envfile ## Re-run the database migrations against the running Postgres
	@$(COMPOSE) --profile '*' run --rm migrate


.PHONY: infra.patroni.list
infra.patroni.list: ## Show the Patroni cluster members and roles (Postgres MODE=ha only)
	@$(COMPOSE) --profile '*' exec -T patroni-1 /opt/patroni/bin/patronictl -c /etc/patroni/patroni.yml list
