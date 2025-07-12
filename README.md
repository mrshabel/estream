# ESTREAM

A simple kafka demo for event-driven communication between microservices.

> [!NOTE]
>
> WSL is required to run the services on windows.

## Setup

1. Start the Kafka instance using docker (No zookeeper required as the instance utilizes KRaft)

```bash
	docker compose up -d
```

## Usage

Run all the service binaries

```bash
    cd bin
    # run the order service
    ./order
    # run the inventory service
    ./inventory
    # run the notification service
    ./notification
```

## Considerations

-   The unofficial Kafka docker image from bitnami is used since it comes with KRaft by default, eliminating the need for a separate zookeeper server.
-   All microservice implementations are mocked for demonstration purposes
-   Naming Conventions used:
    -   Client ID: `<service-name>-service`. eg: `order-service`
    -   Topic: `<entity>.<event-type>`. eg: `order.placed`, `inventory.reserved`
    -   Consumer Group: `<client-id>-group`. eg: `order-service-group`
