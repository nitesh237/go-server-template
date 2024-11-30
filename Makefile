# great stuff here for `make help` @ https://gist.github.com/prwhite/8168133
# COLORS for help sections
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

# global variables for the receipes
DEFAULT_PGDB_USER=root
DEFAULT_PGDB_DATABASE=sample_db
DB_NAME=${DEFAULT_PGDB_DATABASE}
PGDB_ENDPOINT=localhost
PGDB_PORT=5432
MIGRATION_URL="postgres://${DEFAULT_PG_USER}:@${PGDB_ENDPOINT}:${PGDB_PORT}"


create-migration: 
	migrate create -ext sql -dir ./db/${DB_NAME}/migrations -seq ''

_create-db:
	psql -d ${DEFAULT_PGDB_DATABASE} -U ${DEFAULT_PGDB_USER} -tc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" | grep -q 1 || psql -d ${DEFAULT_PGDB_DATABASE} -U ${DEFAULT_PGDB_USER} -c "CREATE DATABASE ${DB_NAME}"

## Take latest db snapshot in latest.sql
_snapshot-db:
	@echo '> Setting up latest.sql'
	@$(eval VERSION = $(shell migrate -database ${MIGRATION_URL}/${DB_NAME}?sslmode=disable -path ./db/${DB_NAME}/migrations version 2>&1 | cat))
	@{ echo '--> Migration Version: $(VERSION) \n' & \
	pg_dump --schema-only --no-owner --no-privileges --no-security-labels --no-tablespaces ${DB_NAME};} | sed '/^SET/d' | sed '/^SELECT/d' | grep -v "^--" | grep "\S"> db/${DB_NAME}/latest.sql

_upgrade-db: _create-db
	migrate -database "${MIGRATION_URL}/${DB_NAME}?sslmode=disable" -path ./db/${DB_NAME}/migrations up ${N}

upgrade-db:	_create-db _upgrade-db _snapshot-db

downgrade-db: 
	migrate -database "${MIGRATION_URL}/${DB_NAME}?sslmode=disable" -path ./db/${DB_NAME}/migrations down ${N}

