#!/bin/bash
go env -w GOOS=linux
go run ./cli/main.go build
