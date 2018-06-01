#!/bin/sh

 echo "Building admin service..."
 cd admin
go build -o ../admin.bin
 cd ..

echo "Building apigateway service..."
 cd apigateway
go build -o ../apigateway.bin
 cd ..

echo "Building auth service..."
 cd auth
go build -o ../auth.bin
 cd ..

echo "Building events service..."
 cd events
go build -o ../events.bin
 cd ..

echo "Building logger service..."
 cd logger
go build -o ../logger.bin
 cd ..

echo "Building registry service..."
 cd registry
go build -o ../registry.bin
 cd ..

echo "Building todo service..."
 cd todo
go build -o ../todo.bin
 cd ..
