#!/bin/bash
go get github.com/githubnemo/CompileDaemon
CompileDaemon -command="go run main.go" -exclude-dir=[.git]
