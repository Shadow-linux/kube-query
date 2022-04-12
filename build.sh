#!/bin/bash

PROJECT_NAME="kubectl-query"

rm ./bin/${PROJECT_NAME}
go build -o ./bin/${PROJECT_NAME} ./main.go
chmod a+x ./bin/${PROJECT_NAME}