#!/bin/sh

set -e
(
    cd "$(dirname "$0")"
    go build -o /tmp/goInterpreter app/*.go
)

exec /tmp/goInterpreter "$@"