PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)
GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

.PHONY: .install-linter
install-linter:#
	# INSTALL GOLANGCI-LINT ###
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.62.2

#lint: .install-linter
.PHONY: lint
lint:
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yaml --verbose

#lint-fast: .install-linter
.PHONY: lint-fast
lint-fast: 
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yaml

sort:
	goimports -w .

swagger:
	@echo "Проверка установки swag..."
	@if ! command -v swag &> /dev/null; then \
		echo "Установка swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@echo "Генерация Swagger документации..."
	$(shell go env GOPATH)/bin/swag init -g app/main.go -o api/swagger/api

test:
	go test ./...

cover:
	go test ./... -coverprofile=coverage.txt -covermod

swagger-update:
	cd api/swagger/api
	git pull

codegen:
# make swagger-update
#	oapi-codegen --config=./api/gen_configs/appinfo/server_cfg.yaml -o ./internal/router/http/appinfo/app_info.gen.go ./api/swagger/api/backend/app_info.yaml

	@echo "Проверка установки OpenAPI Generator..."
	@if ! command -v openapi-generator-cli &> /dev/null; then \
		echo "Установка OpenAPI Generator..."; \
		if command -v npm &> /dev/null; then \
			npm install @openapitools/openapi-generator-cli -g; \
		else \
			echo "Ошибка: npm не установлен. Установите Node.js и npm"; \
			exit 1; \
		fi; \
	fi
	@echo "Генерация Flutter клиента..."
	openapi-generator-cli generate \
		-i api/swagger/api/swagger.yaml \
		-g dart-dio \
		-o clients/flutter \
		--additional-properties=pubName=go_shop_api,pubVersion=1.0.0,pubDescription="Go Shop API Client"

# Генерация всех клиентов
generate-clients: flutter-client
	@echo "Все клиенты сгенерированы"

gqlgen:
	go run github.com/99designs/gqlgen generate


setup-githooks:
	touch .githooks/pre-push
	chmod +x .githooks/pre-push  # Make it executable
	git config core.hooksPath .githooks
	echo "Git hooks directory set to .githooks"