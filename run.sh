#!/bin/sh

echo Testing weaviate...

# Prepare the env
echo Launch database
docker compose up weaviate -d
sleep 10

# Run
echo Run...
go run ./src
