@echo off
go env -w GOOS=windows
go run ./cli/main.go build -v
.\aliax.exe dev-deploy --verbose
del .\aliax.exe
