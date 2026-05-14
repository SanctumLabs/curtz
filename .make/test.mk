
.PHONY: test
test: ## Runs all tests
	@echo "${YELLOW} Running tests ${NC}"
	go test ./...
	@echo "${YELLOW} Done Running tests ${NC}"

.PHONY: test.coverage
test.coverage: ## Runs all tests with coverage
	@echo "${GREEN} Running tests with coverage ${NC}"
	go test -tags testing -v -cover -covermode=atomic -coverprofile=coverage.out ./...
	@echo "${GREEN} Done running tests with coverage ${NC}"

.PHONY: generate.mocks
generate.mocks: ## generates mocks usage: make generate.mocks MOCK_SOURCE=internal/infra/database/postgresql/queries.go MOCK_DESTINATION=internal/infra/database/postgresql/mocks/queries_mock.go MOCK_PACKAGE=postgresqlmocks
	@echo "${GREEN} >>>> Generating mocks for $(MOCK_SOURCE) in $(MOCK_DESTINATION) ${NC}"
	mockgen -destination=$(MOCK_DESTINATION) -package=$(MOCK_PACKAGE) -source=$(MOCK_SOURCE)
	@echo "${GREEN} >>>> Done generating mocks for $(MOCK_SOURCE) in $(MOCK_DESTINATION) ${NC}"
