#!/usr/bin/env bash

echo
echo "==> Running go test <=="
if [ $(uname) == "Linux" ]; then
  go test -msan ./... || exit_code=1
fi
go test -cover -race -coverprofile=cover.out -outputdir=coverage ./... || exit_code=1
echo
echo "==> Running coverage <=="
go tool cover -func=./coverage/cover.out || exit_code=1

exit $exit_code
