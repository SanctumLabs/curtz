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
setup: setup-linting setup-trivy install install-git-hooks

# Will setup linting tools
setup-linting:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.46.2
	chmod +x ./bin/golangci-lint

setup-trivy:
	curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b ./bin v0.16.0

install-git-hooks:
	@echo "Installing Git hooks..."
	cp ./githooks/prepare-commit-msg.sh .git/hooks/prepare-commit-msg
	cp ./githooks/pre-commit.sh .git/hooks/pre-commit
	cp ./githooks/pre-push.sh .git/hooks/pre-push

# Download dependencies
install:
	go mod download

# remove unused dependencies
tidy:
	go mod tidy -v

# Runs project
run:
	go run app/cmd/main.go

run-dev:
	air

# run dev in debug mode
run-dev-d:
	air -d

test:
	go test ./...

test-coverage:
	go test -tags testing -v -cover -covermode=atomic -coverprofile=coverage.out ./...

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
	docker run --rm -i -v hadolint.yaml:/.config/hadolint.yaml hadolint/hadolint < Dockerfile	

lint:
	./bin/golangci-lint run ./...

build:
	@echo "Building application"
	go build -o $(BIN_DIR)/curtz app/cmd/main.go

all: install lint lint-docker test