# Default variables
DB_NAME ?= db_user_service
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= user_service_admin
DB_PASSWORD ?= VCNNGJjsLSmoU5nxnSBBs
MIGRATION_STEPS ?= 1

# https://stackoverflow.com/questions/6405127/how-do-i-specify-a-password-to-psql-non-interactively
connect:
	psql postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}

.PHONY: migrate
migrate:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME} up
	# Target allowing to migrate the db to latest defined step.

migrateOne:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up ${MIGRATION_STEPS}

demigrate:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down

demigrateOne:
	migrate -path ${PWD}/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down ${MIGRATION_STEPS}
