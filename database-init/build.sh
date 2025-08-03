#!/bin/bash
go mod download && go build -o db-init . && ./db-init