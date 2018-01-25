#!/bin/sh

 echo "Compiling protocol..."
 ./proto_gen.sh ./proto

 echo "Building admin service..."
 cd admin
go build -o ../build/admin.bin
 cd ..

echo "Building apigateway service..."
 cd apigateway
go build -o ../build/apigateway.bin
 cd ..

echo "Building auth service..."
 cd auth
go build -o ../build/auth.bin
 cd ..

echo "Building events service..."
 cd events
go build -o ../build/events.bin
 cd ..

echo "Building logger service..."
 cd logger
go build -o ../build/logger.bin
 cd ..

echo "Building registry service..."
 cd registry
go build -o ../build/registry.bin
 cd ..

echo "Building todo service..."
 cd todo
go build -o ../build/todo.bin
 cd ..
