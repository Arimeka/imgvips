#!/usr/bin/env bash

echo
echo "==> Running go vet <=="
go vet ./... || exit_code=1
echo
echo "==> Running golangci-lint <=="
golangci-lint run || exit_code=1

exit $exit_code
