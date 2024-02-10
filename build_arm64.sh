#!/bin/sh

cd ./src
GOOS=darwin GOARCH=arm64 go build -o ./../build/arm64/XlsToRune.arm64
