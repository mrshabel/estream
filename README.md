# ESTREAM

A simple kafka demo for event-driven communication between microservices.

## Setup

1. Start the Kafka instance using docker (No zookeeper required as the instance utilizes KRaft)

```bash
	docker compose up -d
```

## Usage

## Considerations

-   The unofficial Kafka docker image from bitnami is used since it comes with KRaft by default, eliminating the need for a separate zookeeper server.
-   All microservice implementations are mocked for demonstration purposes
