#!/bin/bash

go test $(go list ./... | grep -vE "(vendor)|(test$)|(array$)|(mocked)") -race -coverprofile=coverage.out