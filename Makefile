DOCKER_COMPOSE_FILE := ./docker/docker-compose.yaml
.PHONY: start compile

start: # Start kafka cluster
	docker compose -f $(DOCKER_COMPOSE_FILE) up

compile: # Compile all services into bin directory
	@echo "Compiling all services into 'bin' directory"
	go build -o bin/order cmd/order/*.go
	go build -o bin/inventory cmd/inventory/*.go
	go build -o bin/notification cmd/notification/*.go