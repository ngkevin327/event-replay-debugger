# Contributing

## Branching

- `main` — stable integration branch
- `feat/<area>-<short-description>` — features
- `fix/<area>-<short-description>` — bug fixes
- `chore/<area>-<short-description>` — tooling and maintenance

## Commits

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
type(scope): short description
```

Common types: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`, `perf`.

Keep commits focused; prefer multiple small commits over one large dump.

## Pull requests

- Link related incident or design context when applicable
- Ensure CI passes (`make lint`, `make test`)
- Request review from CODEOWNERS for touched paths
- Update OpenAPI or schema when changing public contracts

## Code style

- Go: `golangci-lint` via `make lint`
- Web: ESLint + Prettier via pnpm in `apps/web`
- Run `make lint` before pushing

## Local setup

See [docs/local-dev.md](docs/local-dev.md).
