#!/usr/bin/env bash

echo
echo "==> Running go vet <=="
go vet ./... || exit_code=1
echo
echo "==> Running golint <=="
golint -set_exit_status ./... || exit_code=1
echo
echo "==> Running gocritic <=="
gocritic check ./... || exit_code=1
echo
echo "==> Running ineffassign <=="
ineffassign ./* || exit_code=1
echo
echo "==> Running gocyclo <=="
gocyclo -over 10 . || exit_code=1

exit $exit_code
