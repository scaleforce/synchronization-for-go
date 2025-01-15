## Synchronization SDK for Go

### Overview

### Quick Start

#### Pub 

Example app that publishes an event every 10 seconds.

1. Add config.env to cmd/pub

```env
# cmd/pub/config.env

AZURE_SERVICEBUS_CONNECTION_STRING=
AZURE_SERVICEBUS_NAMESPACE=
AZURE_SERVICEBUS_TOPIC=
```

2. Run

```shell
go run cmd/pub/main.go
```

#### Sub

Example app that subscribes to an Azure ServiceBus topic to fetch the published events every 10 secondes.

1. Add config.env to cmd/sub

```env
# cmd/sub/config.env

AZURE_SERVICEBUS_CONNECTION_STRING=
AZURE_SERVICEBUS_NAMESPACE=
AZURE_SERVICEBUS_TOPIC=
AZURE_SERVICEBUS_SUBSCRIPTION=
```

2. Run

```shell
go run cmd/sub/main.go
```
