#!/bin/bash
go env -w GOOS=linux
go run ./cli/main.go build --verbose
./aliax dev-deploy --verbose
rm ./aliax
