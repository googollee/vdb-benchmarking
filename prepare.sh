#!/bin/sh

mkdir -p data

echo Generate vectors to data/sources.json...
go run ./src/cli/generate -kline 10 -output data/sources.json

echo Generate queries to data/queries.json...
go run ./src/cli/generate -kline 1 -output data/queries.json
