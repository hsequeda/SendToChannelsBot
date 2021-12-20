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
	sqlboiler -d --wipe -o adapter/model/ -p model psql --add-global-variants --no-hooks

integration_test:
	go test --race -v -count=1 -p=8 -parallel=8 ./...

.PHONY: run
run:
	go run main.go
