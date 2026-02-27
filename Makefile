.PHONY: dev dev-deps setup-local verify-local migrate-local build test lint clean

GO_SERVICES := $(shell find apps services packages -name go.mod -exec dirname {} \; 2>/dev/null | sort -u)
COMPOSE := docker compose -f infra/docker-compose/docker-compose.yml
DATABASE_URL ?= postgres://replay:replay@localhost:15433/replay?sslmode=disable

dev-deps:
	$(COMPOSE) up -d

setup-local: dev-deps migrate-local
	@test -f .env.local || cp .env.local.example .env.local
	pnpm install --filter @replay/web
	@echo "Setup complete. See docs/local-dev.md"

migrate-local:
	cd apps/api-gateway && DATABASE_URL="$(DATABASE_URL)" go run ./cmd/migrate -dir ./migrations

verify-local:
	@./scripts/verify-local.sh

dev: dev-deps
	@echo "Dependencies up. Run: make setup-local && see docs/local-dev.md"

build:
	@for dir in $(GO_SERVICES); do \
		echo "building $$dir..."; \
		(cd $$dir && go build ./...) || exit 1; \
	done

test:
	@for dir in $(GO_SERVICES); do \
		(cd $$dir && go test ./...) || exit 1; \
	done
	@if [ -f apps/web/package.json ]; then pnpm --filter @replay/web test 2>/dev/null || true; fi

lint:
	@./scripts/lint-go.sh
	@if [ -f apps/web/package.json ]; then pnpm --filter @replay/web lint 2>/dev/null || true; fi

clean:
	rm -rf bin/ dist/
	find . -name '*.test' -delete 2>/dev/null || true
