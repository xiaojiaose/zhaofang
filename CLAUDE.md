# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Dev Commands

```bash
# Backend dev
cd server && go run .

# Frontend dev
cd web && npm install && npm run serve

# Frontend production build
cd web && npm run build

# Backend build verification
cd server && go build ./...

# Swagger regen (after changing API annotations)
cd server && swag init

# Full packaged build (from repo root)
make build
```

## Architecture

### Layer mapping
- **API layer**: `server/api/v1/<package>/` — HTTP handlers, one file per domain
- **Service layer**: `server/service/<package>/` — business logic, may be split across multiple files for the same service
- **Model layer**: `server/model/<package>/` — DB entities; `server/model/system/request/` and `server/model/system/response/` for DTOs
- **Router layer**: `server/router/<package>/` — route registration, delegates to api layer

### Key structural patterns
- **Multiple API packages**: `api/v1/system/`, `api/v1/resource/`, `api/v1/center/` — each registers into a separate `ApiGroup` sub-struct in `api/v1/enter.go`
- **Service group assembly**: `service/<package>/enter.go` declares the `ServiceGroup` struct; individual service structs are defined in separate files in the same directory
- **Router wiring**: `server/initialize/router.go` calls `Init*Router()` for each group; each router file (e.g., `router/system/sys_xiaoqu.go`) defines its own route registrations
- **Domain split**: `model/house/` contains `dict_building`, `dict_unit`, `dict_house` entities (楼盘字典); `model/system/` contains `sys_xiaoqu` (小区). These are separate model packages despite being related domain-wise.

### Domain: XiaoQu and 楼盘字典
The `xiaoqu` (小区) and `dict_building/dict_unit/dict_house` (楼盘字典) modules share a logical relationship:
- `xiao_qu` → `dict_building` via `community_id`
- `dict_building` → `dict_unit` via `building_open_id`
- `dict_unit` → `dict_house` via `unit_open_id`

API handlers are in `api/v1/system/sys_xiaoqu.go` (main xiaoqu routes) and `api/v1/resource/sys_xiaoqu.go` (base/联动查询 routes). **Both files define a `XiaoQuApi` struct — they are different types in different packages.**

The service logic is split:
- `service/system/sys_xiao_qu.go` — `XiaoQuService` handles Create/Edit/List/GetInfo and auto-fills `community_id`
- `service/system/sys_xiaoqu_dict.go` — `XiaoQuService` handles `GetDictTree`/`UpsertDict`/`DeleteDict`

### Frontend API layer
- `web/src/utils/request.js` — axios instance with auth headers (`x-token`, `x-user-id`)
- Dynamic routes loaded after login via Pinia store; `web/src/router/index.js` only declares static shell routes
