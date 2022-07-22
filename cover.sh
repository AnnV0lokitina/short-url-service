#!/bin/bash

go test $(go list ./... | grep -vE "(vendor)|(test)|(array$)|(mocked)|(bench)|(checker)|(proto)") -race -count=1 -coverprofile=coverage.out -args local
go tool cover -func=coverage.out