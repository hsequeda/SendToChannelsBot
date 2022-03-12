include .env
export

# database
.PHONY: migration_new
migration_new:
	sql-migrate new -config=sql/dbconfig.yml $(name)

.PHONY: migration_run
migration_run:
	sql-migrate up -config=sql/dbconfig.yml

.PHONY: migration_revert
migration_revert:
	sql-migrate down -config=sql/dbconfig.yml

.PHONY: migration_status
migration_status:
	sql-migrate status -config=sql/dbconfig.yml

.PHONY: generate_db_model
generate_db_model:
	sqlboiler -d --wipe -o pkgs/common/sql/psql_model/ -p psql_model psql --add-global-variants --no-hooks

.PHONY: test
integration_test:
	@echo "Start integration tests...."
	@./scripts/integration_test.sh .integration_test.env

.PHONY: run
run:
	go run main.go

.PHONY: generate_openapi_server
generate_openapi_server:
	oapi-codegen -generate types -o pkgs/account/port/rest/openapi_types.gen.go -package rest api/open_api/account.yml
	oapi-codegen -generate chi-server -o pkgs/account/port/rest/openapi_api.gen.go -package rest api/open_api/account.yml
	oapi-codegen -generate types -o pkgs/administration/port/rest/openapi_types.gen.go -package rest api/open_api/administration.yaml
	oapi-codegen -generate chi-server -o pkgs/administration/port/rest/openapi_api.gen.go -package rest api/open_api/administration.yaml
