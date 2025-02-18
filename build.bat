@echo off
go env -w GOOS=windows
go run ./cli/main.go build --verbose
.\aliax.exe dev-deploy --verbose
del .\aliax.exe
