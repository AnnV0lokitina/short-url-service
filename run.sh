#!/bin/bash

go run -tags profile -ldflags "-X main.buildVersion=v1.0.1 \
                -X 'main.buildDate=$(date +'%Y/%m/%d')' \
                -X 'main.buildCommit=My comment'" \
                github.com/AnnV0lokitina/short-url-service/cmd/shortener