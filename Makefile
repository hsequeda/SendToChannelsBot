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
	rm -rf internal/common/sql/psql_models/*
	sqlboiler -o internal/common/sql/psql_models -p psql_models psql --add-global-variants --no-hooks --no-auto-timestamps
