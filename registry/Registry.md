# Registry

Registry service keeps a list of all service types and service type instances, performs health checks and performs load balancing.

## Environment variables for configuration

- `RG_RPC_PORT` - port for RPC server
- `RG_HTTP_PORT` - port for HTTP server
- `RG_HTTPS_PORT` - port for HTTPS server. If `RG_HTTPS_PORT` is not set or set to 0, then HTTP server is used
- `RG_TIMEOUT` - number of seconds for network operations timeout
- `RG_RETRIES` - number of retries to perform for network operations

## REST API

- `GET /registry/health`:
    - Information:
    
        Returns health status of the Registry service. Good health requires Registry service to have RPC and 
        HTTP(w/ or w/o TLS) servers running.
    - Returns:
        ```json
        {
          "up": true
        }
        ```
        
        Values for `up`:
        - `true` - service is running RPC and HTTP servers
        - `false` - service is not running RPC or HTTP server, so some functions are unavailable
- `GET /registry/service/types`:
    - Information:
        
        Returns list of registered service types.
    - Returns:
        ```json
        {
          "types": [
            {
              "type": "<type>"
            },
            {
              "type": "<type>"
            }
          ]
        }
        ```
        
- `GET /registry/service/types/{type}`:
    - Information:
        
        Returns information about all registered instances of `type`
    - Returns:
        ```json
        {
          "services": [
            {
              "proto": "<proto>",
              "type": "<type>",
              "host": "<host>",
              "port": "<port>",
              "httpPort": "<httpPort>",
              "httpsPort": "<httpsPort>",
              "routes": ["<route>", "<route>"],
              "health": "<health>",
              "weight": "<weight>",
              "signature": "<signature>"
            }
          ]
        }
        ```
        
        Values:
        - `proto`: `http`/`https`
        - `type`: string name of service type
        - `host`: ip address
        - `port`: 0-65535
        - `httpPort`: 0-65535
        - `httpsPort`: 0-65535
        - `routes`: array of service endpoints without protocol, host and port
        - `health`: service health check endpoint without protocol, host and port. E.g. `/registry/health`
        - `signature`: JWT token

- `GET /registry/service/types/{type}/best`:
    - Information:
        
        Returns best suitable instance of service `type` by performing load balancing
    - Returns:
        ```json
        {
          "proto": "<proto>",
          "type": "<type>",
          "host": "<host>",
          "port": "<port>",
          "httpPort": "<httpPort>",
          "httpsPort": "<httpsPort>",
          "routes": ["<route>", "<route>"],
          "health": "<health>",
          "weight": "<weight>",
          "signature": "<signature>"
        }
        ```
        
        Values:
        - `proto`: `http`/`https`
        - `type`: string name of service type
        - `host`: ip address
        - `port`: 0-65535
        - `httpPort`: 0-65535
        - `httpsPort`: 0-65535
        - `routes`: array of service endpoints without protocol, host and port
        - `health`: service health check endpoint without protocol, host and port. E.g. `/registry/health`
        - `signature`: JWT token
            
- `GET /registry/service/{instanceName}`:
    - Information:
    
        Returns information about service defined as `instanceName`. `instanceName` is a combination of service type and its index with
        dash in between. E.g. `auth-2` is a third registered instance of `Auth` service.
    - Returns:
        
        Returned JSON is the same as in `GET /registry/service/types/{type}/best`
        
- `GET /registry/service/types/{type}/{index}`:
    - Information:
    
        Basically the same as `GET /registry/service/{instanceName}`, but service type and index are separated.
    - Returns:
    
        Returned JSON is the same as in `GET /registry/service/{instanceName}`
        
- `POST /registry/service`:
    - Information:
    
        Registers service in the registry. The registry performs health check before adding service to its list and checks
        if similar service already exists by comparing existing services' host, port and signature
        
    - Input:
        - Headers:
            - `Content-Type:application/json`
        - Body:
            ```json
            {
              "proto": "<proto>",
              "type": "<type>",
              "host": "<host>",
              "port": "<port>",
              "httpPort": "<httpPort>",
              "httpsPort": "<httpsPort>",
              "routes": ["<route>", "<route>"],
              "health": "<health>",
              "weight": "<weight>",
              "signature": "<signature>"
            }
            ```
            
            Values:
            - `proto`: `http`/`https`
            - `type`: string name of service type
            - `host`: ip address
            - `port`: 0-65535
            - `httpPort`: 0-65535
            - `httpsPort`: 0-65535
            - `routes`: array of service endpoints without protocol, host and port
            - `health`: service health check endpoint without protocol, host and port. E.g. `/registry/health`
            - `signature`: JWT token
    - Returns:
        ```json
        {
          "status": "OK",
          "message": "OK"
        }
        ```
        
        Values:
        - `status`:
            - `OK` - service is registered
            - `INVALID` - signature is invalid
            - `FAIL` - internal error
            - `CANCELED` - server canceled registering
            - `NOT_IMPLEMENTED` - health route is not responding or responded with `"up": false`
            - `NULL` - no service is received
            - `EXISTS` - service is already registered
        - `message` - error description or `"OK"`