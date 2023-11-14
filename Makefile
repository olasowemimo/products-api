include .env

stop_containers:
	@echo "stopping other docker containers..."
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers."; \
		docker stop $$(docker ps -q); \
	else \
		echo "no running containers found."; \
	fi

create_db_container:
	@echo "creating docker database container..."
		docker run --name ${DB_CONTAINER} -p ${DB_PORT}:${DB_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:12-alpine

create_db:
	@echo "creating database..."
	docker exec -it ${DB_CONTAINER} bash -c "createdb --username ${DB_USER} --owner ${DB_USER} ${DB_NAME}"

start_db_container:
	@echo "starting docker database container..."
	docker start ${DB_CONTAINER}

create_db_migration:
	@echo "creating sqlx database migration..."
	sqlx migrate add -r init

migrate_db_up:
	@echo "sqlx database migration up."
	sqlx migrate run --database-url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

migrate_db_down:
	@echo "sqlx database migration down."
	sqlx migrate revert --database-url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

build:
	if [ -f "${BINARY}" ]; then \
		rm ${BINARY}; \
		echo "deleted go ${BINARY}."; \
	fi
	@echo "building go alpine..."
	go build -o ${BINARY} -v ./cmd/api/*.go

run: build
	./${BINARY}

stop:
	@echo "stopping go api server..."
	@-pkill -SIGTERM -f ./${BINARY}
	@echo "go api server stopped."