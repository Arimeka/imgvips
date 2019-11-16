#!/usr/bin/env bash

echo
echo "==> Running go test <=="
go test -cover -race -coverprofile=cover.out -outputdir=coverage ./... || exit_code=1
echo
echo "==> Running coverage <=="
go tool cover -func=./coverage/cover.out || exit_code=1

exit $exit_code
