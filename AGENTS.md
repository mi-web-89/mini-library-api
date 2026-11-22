# AGENTS.md

Go 1.27 + Fiber v3 REST API (`mini-library-api`). No Makefile, no linter config: use `gofmt` + `go vet`.

## Commands

```sh
go run cmd/seed/role.go        # REQUIRED first: creates PostgreSQL enums role + rent_status
go run cmd/api/main.go         # server; AutoMigrate runs on every start
air                            # hot reload (builds ./cmd/api/main.go)
go run cmd/generate/admin.go   # creates admin user from ADMIN_* env vars

go test ./...
go test ./auth/ -run TestLogin_Success -v
```

Run everything from the repo root: `utils.GetConfig` loads `.env` relative to CWD.

## Known broken state (verify before blaming your change)

- `go build ./...` fails: `api/handlers/categories.go` is an empty (0-byte) WIP file.
- `go test ./auth/` fails 3 register tests: tests expect 12 INSERT args, model now has 13 (`profile_picture_id`).

## Config quirk

`APP_MODE != "production"` → reads `.env` via Viper; `APP_MODE=production` → reads only OS env vars. Don't set `APP_MODE=production` locally without exporting the vars.

## Architecture

- Wiring is manual in `api.NewApp` (composition root). No DI framework.
- Handlers take `fiber.Ctx` **by value** (Fiber v3), not `*fiber.Ctx`.
- Each module operation is a separate file with its own struct (`register`, `get`, `update`, ...) embedded into `service` in `service.go`. Add an operation = new file + embed it + add to the `Service` interface.
- Persistence models live in `db/models` (canonical). Request/response structs live in the module (`auth/auth.go`, `modules/users/users.go`, `modules/categories/category.go`). **`User` is duplicated in all three** — change all copies or consolidate.
- GORM is injected as `*gorm.DB` directly into services; there are no repository interfaces. Mocks exist only for `Service` interfaces.
- Responses use `pkg/dtos.Response[T]` (`status`, `message`, `data`); map DB errors with `utils.GetErrorMessage`.
- No migration files: add new models to the `AutoMigrate` list in `db/drivers/postgres.go`.

## Gotchas

- `middlewares.ValidateBody` stores the DTO in `c.Locals("validatedBody")`, but **no handler reads it** — handlers re-bind with `ctx.Bind()` and validate via `utils.ValidateStruct`. Follow the handler pattern.
- `utils.GenerateJWT` is dead code; live JWT logic is `middlewares.JWTConfig`. Validation rules are duplicated in `middlewares/validation.go` and `utils/validation.go`.
- Pagination (`utils.Paginate`) interpolates `Sort` into SQL and relies on `Keyword` being set; don't pass unsanitized user input.
- `modules/categories` implements full CRUD but has no handler/routes yet. `Book` model exists with no module. `pkg/ai`, `pkg/rajaongkir`, and Cloudinary upload are wired but unused by routes.
- `docker-compose.yml` has no `build:` key (expects a prebuilt image) and maps `${PORT}:1323` while the app listens on `${PORT}`.
- `.mockery.yml` still references the old `mini-library-api/users` path (now `modules/users`); fix it before regenerating mocks.
