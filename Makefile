.PHONY: dev dev-deps build test lint clean

GO_SERVICES := $(shell find apps services packages -name go.mod -exec dirname {} \; 2>/dev/null | sort -u)

dev-deps:
	docker compose -f infra/docker-compose/docker-compose.yml up -d

dev: dev-deps
	@echo "Dependencies up. Run services individually from apps/*/cmd."

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
