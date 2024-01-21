#!/bin/sh

echo Generate vectors to data/sources.json...
mkdir -p data
go run ./src/cli/generate -kline 10 -output data/sources.json
