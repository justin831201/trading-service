# Trading Service
## Architecture
The service can separate into two microservices, and they communicate through message queue.

The first microservice called `order-service`.
It provides API interface for users to send orders, and it would produce these orders to message queue. 

The second microservice called `exchange-service`
It receives orders from message queue, and then it tries to match those orders for each specific duration. 
### Why we need message queue
* When the trading volume grows bigger and bigger, we need to scale up `order-service`. Message queue can help us collect orders from different `order-service` and we can get all orders from the queue.
* It can treat as buffer place when the requests are too many and `order-service` can not consume these requests immediately.

## Interface
We provide an endpoint for you to send order.
```
path: /order
method: POST
request:
    json:
        user_id: string
        order_type: 0 or 1  // 0 means buy, 1 means sell
        quantity: int
        price_type: 0 or 1  // 0 means market price, 1 means limit price
        price: int
response:
    json:
        order_id: string
```

cURL example
```shell
$ curl -X POST -d '{"user_id": "9368d4af-ae72-4cc8-931d-ce2268ce8ffb", "order_type": 0, "quantity": 6, "price_type": 1, "price": 100}' http://localhost:8080/order

{"order_id":"a540574c-63be-4f13-8df5-8e7a8347c497"}
```

## Installation
The service can run at local and in docker. You can choose by your development habit.

### Local
#### Prerequisite
* kafka
#### Config Settings
you need to set kafka brokers URI in your configs that microservices can connect to the brokers.
The config paths are `config/exchange-service/config-local.yaml` and `config/order-service/config-local.yaml`
```yaml
...
kafka:
  bootstrap_servers:
    - ${YOUR_BOOTSTRAP_SERVERS}
...
```
#### Run
you can run the service by the following commands.
```shell
(Terminal 1) Run exchange service
$ go run cmd/exchange-service/main.go -c config/exchange-service/config-local.yaml

(Terminal 2) Run order service
$ go run cmd/order-service/main.go -c config/order-service/config-local.yaml
```

### Docker
#### Prerequisite
* docker
#### Run
You can run the service by the following command.
```shell
$ docker-compose up --build
```

## Experiment
You can use `/script/trading-client.sh` to generate some random orders.
Ant Then you can watch `exchange-service` logs to ensure the algorithm are correct.

```shell
$ ./script/trading-client.sh
```

## TODO List
* Write OpenAPI Document
* Implement graceful Shutdown
* Implement notification microservice to send trading notification to user
* Write unit and integration test cases
* Implement CI/CD process
