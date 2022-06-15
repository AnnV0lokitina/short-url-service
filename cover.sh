#!/bin/bash

go test $(go list ./... | grep -vE "(vendor)|(test)|(array$)|(mocked)|(bench)") -race -coverprofile=coverage.out
go tool cover -func=coverage.out