#!/bin/bash

set -e

# CONFIGURATIONS
WORK_DIR=.cover
COVERAGE_PROFILE_FILE="$WORK_DIR/cover.out"
COVERAGE_MODE=count

go get github.com/stretchr/testify

generateCoverage() {
    rm -rf "$WORK_DIR"
    mkdir "$WORK_DIR"

    for item in "$@"; do
        f="$WORK_DIR/$(echo $item | tr / -).cover"
        go test -covermode="$COVERAGE_MODE" -coverprofile="$f" "$item"
    done

    echo "mode: $COVERAGE_MODE" >"$COVERAGE_PROFILE_FILE"
    grep -h -v "^mode:" "$WORK_DIR"/*.cover >>"$COVERAGE_PROFILE_FILE"
}

showCoverageReport() {
    go tool cover -${1}="$COVERAGE_PROFILE_FILE"
}

if [ "$1" == "--jenkins" ]; then
	go test -v ./... | go-junit-report > ./junit.xml
else
    generateCoverage $(go list ./...)
    showCoverageReport func

    if [ "$1" == "--html" ]; then
        showCoverageReport html
    fi
fi