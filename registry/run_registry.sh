#!/bin/bash

docker run -d -it -p 127.0.0.1:3001:3000 \
-p 127.0.0.1:8080:80 \
-e RG_HTTP_PORT='80' \
-e RG_RPC_PORT='3000' \
-e RG_TIMEOUT='30' \
micro-todo
