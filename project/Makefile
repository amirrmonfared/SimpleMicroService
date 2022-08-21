FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailerApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"


## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mailer build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"
## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker compose down
	@echo Done!

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo Building logger binary...
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo Done!

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo Building logger binary...
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo Done!

## build_mailer: builds the mailer binary as a linux executable
build_mailer:
	@echo Building mailer binary...
	cd ../mailer-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo Done!

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo Building broker binary...
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo Done!

## build_auth: builds the authentication binary as a linux executable
build_auth:
	@echo Building authentication binary...
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo Done!

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY}

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"
test:
	@echo Starting Test...
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes
	@echo Done!

.PHONY: up up_build down build_broker build_front start stop build_logger