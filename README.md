# Micro TODO

Trying to write a TODO app with micro services architecture. I am using several technologies just for the sake of using them
and getting first-hand experience with them.

## Installation

1. Install go-1.9
2. Install Docker
3. Install `dep` - [dep](github.com/golang/dep/cmd/dep)
4. Install dependencies by running `dep ensure`
5. Run `docker build -t micro-todo .` to build Docker image
6. Start Docker container from built image
    > By default, container runs `Registry` service. To specify which service to run, append `<service name>.bin` to the 
    container running command. E.g. `docker run -d -p 127.0.0.1:30010:3000 -p 127.0.0.1:8080:80 micro-todo auth.bin`
    
    > To start admin web application, go to `admin/README.md` for instructions

## Project structure

- `admin` - `Admin` web application for exploring services state
- `auth` - `Auth` service
- `logger` - `Logger` service
- `proto` - shared service RPC protocol and generated API gateways
- `registry` - `Registry` service
- `todo` - `TODO` service
- `vendor` - project dependencies

## Services

### Service list

1. `Registry` - **public** - list of all available services and load balancing
2. `Auth` - **internal** - authorizes and registers users
3. `TODO` - **internal** - keeps TODO tasks
4. `Admin` - **admin** - provides WEB UI for logs and status of all services
5. `Logger` - **internal** - receives logs from all services and provides them to `Admin`

### Service API and configuration

For information about service API and configuration see following files:
- [Registry](registry/Registry.md)

### Communication process

#### Clients's communication with `Auth` and `TODO` through `Registry`

0. Services register themselves in `Registry`
1. Client performs request to the `Registry` to specific service and its API endpoint
2. `Registry` proxies request to the best instance, converting request to internal protocol, and waits for response
3. On receiving response, `Registry` converts response from internal protocol to HTTP response and sends back to client

#### `Registry`, `Auth`, `TODO` services' communication with `Logger` service

0. Service queries `Registry` for `Logger` service
1. Service sends log message to `Logger` service

#### `Admin` service's communication with `Logger`

0. `Admin` queries `Registry` for `Logger` service
1. `Admin` queries `Logger` service for requested service type logs
2. `Admin` updates log every `n` seconds
3. `Admin` queries `Registry` for services status  

### Features

- `Registry`:
    - [x] Works on gRPC
    - [x] Works on HTTP
    - [x] Keeps up-to-date list by performing service instance health checks
    - [x] Supports server-side registering
    - [x] Responses with service types list
    - [x] Responses with service type instances list
    - [x] Responses with address of suitable requested service (balancing)
    - [x] Supports several load balancing algorithms:
        - [x] Random balancer
        - [x] Round-robin balancer
        - [x] Weighted random balancer
        - [x] Weighted round-robin balancer
    - [x] Has health HTTP endpoint
    - [x] Has health RPC method
    - [x] Responses with error message if service is unavailable
    - [ ] Restricts querying internal services from external clients
    - [ ] Restricts querying admin services from external clients
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
    - [x] Shows health status of all services
    - [x] Periodically updates status and logs
- `Logger`:
    - [ ] Receives logs from different services
    - [ ] Stores them by service type instances
    - [ ] Responses with logs of specific service type instance
    - [ ] Persists logs on disks
    