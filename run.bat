@echo off
go build main.go taskmgr.go cmd.go server.go
main
del main.exe