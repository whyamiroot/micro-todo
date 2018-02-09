#!/bin/bash

docker build -t todo-registry ./registry

docker run -d -p 127.0.0.1:3001:3000 \
-p 127.0.0.1:8080:8080 \
-v /home/lev/go/src/github.com/whyamiroot/micro-todo:/go/src/github.com/whyamiroot/micro-todo todo-registry \
--build-arg CACHEBUST=$(date +%s)
