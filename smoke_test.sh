#!/usr/bin/env bash
set -e
go mod tidy
go run smoke_check.go
