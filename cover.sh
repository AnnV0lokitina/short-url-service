#!/bin/bash

export LOCAL="local"
go test $(go list ./... | grep -vE "(vendor)|(test)|(array$)|(mocked)|(bench)|(checker)") -race -coverprofile=coverage.out
go tool cover -func=coverage.out