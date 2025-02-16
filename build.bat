@echo off
go env -w GOOS=windows
go run ./cli/main.go build
