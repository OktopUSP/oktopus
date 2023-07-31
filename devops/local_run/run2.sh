#!/bin/bash

cd $HOME/dev/oktopus/backend/services/controller
go run cmd/oktopus/main.go -u root -P root -mongo mongodb://172.16.238.3:27017
echo ""
bash