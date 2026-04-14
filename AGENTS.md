# AGENTS.md

## Repo shape
- Monorepo-style split: `server/` is the Go backend, `web/` is the Vue 3/Vite frontend, root `Makefile` orchestrates cross-project build/package flows.
- Real backend entrypoints: `server/main.go` → `initializeSystem()` → `core.RunServer()`; routing is wired in `server/initialize/router.go`.
- Real frontend entrypoints: `web/src/main.js` and `web/src/router/index.js`.
- Dynamic app routes are **not** declared in `web/src/router/index.js`; they are added after login via the Pinia router store. Only shell routes (`/login`, `/init`, etc.) are static.

## Source-of-truth commands
- Backend dev: `cd server && go run .`
- Frontend dev: `cd web && npm install && npm run serve`
- Frontend production build: `cd web && npm run build`
- Swagger regen: `cd server && swag init` or `make doc`
- Full packaged build: run from repo root with `make build`
- Local split builds: `make build-web-local` and `make build-server-local`
- Docker image flows: `make image` (combined image) or `make images` (web/server/all)

## Verification rules
- Trust executable config over prose: `server/go.mod` currently declares `go 1.23.0` with `toolchain go1.24.3`, even though README/CI still mention Go 1.22.
- Do **not** assume `npm run test` or `npm run lint` exist just because `web/README.md` says so; `web/package.json` only defines `serve`, `build`, `preview`, `limit-build`, and `fix-memory-limit`.
- Practical focused verification is usually:
  - frontend-only change: `cd web && npm run build`
  - backend-only change: `cd server && go build ./...`
  - swagger/config/init changes: also verify `cd server && swag init` if API annotations changed

## Repo-specific gotchas
- `web/package.json` runs `node openDocument.js` before Vite in `npm run serve`; that file is copyright/protection text plus OS-specific `open` logic. Do not remove it casually.
- `Makefile` local web build uses a Tencent Yarn mirror (`yarn config set registry http://mirrors.cloud.tencent.com/npm/`). If installs behave strangely, check whether that registry choice is the cause.
- Backend startup is not lightweight: `initializeSystem()` auto-migrates tables, initializes ZincSearch, and calls `initialize.GeoXiaoqu()`, which loads all rows from `xiao_qu` to build an in-memory geo index.
- `server/initialize/router_biz.go` is currently a placeholder; most real route wiring happens in `server/initialize/router.go` plus the group entry files under `server/router/`, `server/api/v1/`, and `server/service/`.
- Backend exposes MCP SSE/message endpoints from config (`server/initialize/mcp.go`); check `server/config*.yaml` before changing those paths.

## Config and safety
- `server/config.yaml` contains live-looking secrets/credentials. Treat it as sensitive; do not paste values into chat output, and avoid committing edits unless the task explicitly requires config rotation.
- Safer example/default config for containers lives in `server/config.docker.yaml`; use that as the reference when you need placeholder values.
- Frontend API base URL comes from Vite env files and is consumed in `web/vite.config.js` and `web/src/utils/request.js`.
- Auth headers are repo-specific: frontend sends both `x-token` and `x-user-id` in `web/src/utils/request.js`.

## Useful file anchors
- Root build/deploy workflow: `Makefile`, `.github/workflows/ci.yaml`
- Backend structure overview: `server/README.md`, `server/initialize/router.go`, `server/initialize/gorm.go`
- Frontend structure overview: `web/package.json`, `web/vite.config.js`, `web/src/pinia/modules/user.js`
