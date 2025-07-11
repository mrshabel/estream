DOCKER_COMPOSE_FILE := ./docker/docker-compose.yaml
.PHONY: start

start: # Start kafka cluster
	docker compose -f $(DOCKER_COMPOSE_FILE) up