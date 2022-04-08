BIN_DIR = $(PWD)/bin

OSFLAG 				:=
ifeq ($(OS),Windows_NT)
	OSFLAG += -D WIN32
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		OSFLAG += -D AMD64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OSFLAG += -D IA32
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		OSFLAG += -D LINUX
	endif
	ifeq ($(UNAME_S),Darwin)
		OSFLAG += -D OSX
	endif
		UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		OSFLAG += -D AMD64
	endif
		ifneq ($(filter %86,$(UNAME_P)),)
	OSFLAG += -D IA32
		endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		OSFLAG += -D ARM
	endif
endif

echoos:
	@echo $(OSFLAG)

# will setup all the tools required by the project
setup: setup-linting setup-hadolint setup-trivy

# Will setup linting tools
setup-linting:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.41.0
	chmod +x ./bin/golangci-lint

# TODO: setup hadolint based on the OS of the development machine this is being run on
# this can be done with the OSFLAG above which already detects the current OS. Currently, this only setups up
# hadolint on Linux
setup-hadolint:
# if running on Linux
	wget -O ./bin/hadolint https://github.com/hadolint/hadolint/releases/download/v1.16.3/hadolint-Linux-x86_64
	chmod +x ./bin/hadolint

setup-trivy:
	curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b ./bin v0.16.0

# Download dependencies
install:
	go mod download

# remove unused dependencies
tidy:
	go mod tidy -v

# Runs project
run:
	go run cmd/app/main.go

test:
	go test ./...

test-coverage:
	go test -tags testing -v -cover -covermode=atomic ./...

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

clean:
	if [ -f ${BIN_DIR} ] ; then rm ${BIN_DIR} ; fi

scan-docker-image:
	@echo "Scanning Docker Image: $(IMAGE)"
	./bin/trivy $(IMAGE)

# See local hadolint install instructions: https://github.com/hadolint/hadolint
# This is linter for Dockerfiles
lint-docker:
	@echo "Running lint checks on Dockerfile"
	./bin/hadolint Dockerfile

lint:
	./bin/golangci-lint run ./...

build:
	@echo "Building application"
	go build -o $(BIN_DIR) cmd/app/main.go

# See https://circleci.com/docs/2.0/local-cli/#processing-a-config
validate-circleci:
	circleci config validate

# See https://circleci.com/docs/2.0/local-cli/#processing-a-config
process-circleci:
	circleci config process .circleci/config.yml

all: install lint lint-docker test