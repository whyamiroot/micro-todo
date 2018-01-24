# Micro TODO

Trying to write a TODO app with micro services architecture

## Project structure

- `admin` - `Admin` service
- `apigateway` - `APIGateway` service
- `auth` - `Auth` service
- `build` - builds output directory
- `events` - `Events` service
- `logger` - `Logger` service
- `registry` - `Registry` service
- `todo` - `TODO` service
- `vendor` - project dependencies
  > Use [dep](github.com/golang/dep/cmd/dep) to install dependencies

## Services

### Service list

1. `Registry` - **public** - list of all available services and load balancing
2. `APIGateway` - **public** - proxies HTTP requests to internal service instances
3. `Auth` - **internal** - authorizes and registers users
4. `TODO` - **internal** - keeps TODO tasks
5. `Admin` - **admin** - provides WEB UI for logs and status of all services
6. `Logger` - **internal** - receives logs from all services and provides them to `Admin`
7. `Events` - **internal** - sends and receives events for service syncing

### Communication process

#### Clients's communication with `Auth` and `TODO` through `APIGateway`

0. Services register themselves in `Registry`
1. Client performs request to the `Registry` to find the `APIGateway`
2. `Registry` finds suitable `APIGateway` instances and responses with address
3. Client performs request to the `APIGateway`
4. `APIGateway` queries `Registry` for a suitable service instance
5. `APIGateway` proxies request to the instance, converting request to internal protocol, and waits for response
6. On receiving response, `APIGateway` converts response from internal protocol to HTTP response and sends back to client

#### `Auth` and `TODO` services' communication with `Events` service

0. Service queries `Registry` for `Events` service
1. Service subscribes to certain events
2. Service sends sync command to the `Events` service
    > For example - `todo` task with id `id` changed from `value` to `value`
3. `Events` service sends sync command to all subscribers

#### `Registry`, `Auth`, `TODO`, `Events`, `APIGateway` services' communication with `Logger` service

0. Service queries `Registry` for `Logger` service
1. Service sends log message to `Logger` service

#### `Admin` service's communication with `Logger`

0. `Admin` queries `Registry` for `Logger` service
1. `Admin` queries `Logger` service for requested service type logs
2. `Admin` updates log every `n` seconds
3. `Admin` queries `Registry` for services status  

#### `APIGateway` service's communication with `Auth` 

0. Service queries `Registry` for the suitable `Auth` service
1. Service authorizes user by querying `Auth` service
2. Service authenticates some action by querying `Auth` service

### Features

- `Registry`:
    - [ ] Works on both gRPC
    - [ ] Works on HTTP
    - [ ] Keeps up-to-date list by performing service instance health checks
    - [ ] Supports self-register requests
    - [ ] Responses with service types list
    - [ ] Responses with service type instances list
    - [ ] Responses with address of suitable requested service (balancing)
    - [ ] Proxies health check requests to service instances
    - [ ] Responses with error message if service is unavailable
    - [ ] Restricts querying internal services for external clients
    - [ ] Restricts querying admin services for external clients
    - [ ] Writes logs to `Logger`
- `APIGateway`:
    - [ ] Proxies HTTP request to a corresponding instance, converting from HTTP to internal protocol
    - [ ] Queries `Registry` for suitable service instance
    - [ ] Responses with HTTP error if service is unavailable
    - [ ] Manages headers
    - [ ] Writes logs to `Logger`
- `Auth`:
    - [ ] Registers new users
    - [ ] Authorizes users
    - [ ] Provides auth tokens for authorized users
    - [ ] Authenticates users by their auth tokens
- `TODO`:
    - [ ] Stores TODO tasks
    - [ ] Stores TODO lists of user
    - [ ] Allows adding TODO tasks
    - [ ] Allows editing TODO tasks
    - [ ] Allows deleting TODO tasks
    - [ ] Allows sharing TODO tasks with different users
    - [ ] Writes logs to `Logger`
- `Admin`:
    - [ ] Authenticates admin in WEB UI
    - [ ] Collects logs from `Logger`
    - [ ] Shows health status of all services
    - [ ] Periodically updates status and logs
- `Logger`:
    - [ ] Receives logs from different services
    - [ ] Stores them by service type
    - [ ] Stores them by service type instance
    - [ ] Responses with logs of specific service type
    - [ ] Responses with logs of specific service type instance
    - [ ] Persists logs on disks
- `Events`:
    - [ ] Receives subscribe requests
    - [ ] Receives events
    - [ ] Notifies subscribers
    - [ ] Writes logs to `Logger`
    