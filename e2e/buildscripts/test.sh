#!/usr/bin/env bash
set -e

if [ -z "${CTLNAME}" ]; 
then
    CTLNAME="m-e2e"
fi

# Create a temp dir and clean it up on exit
TEMPDIR=`mktemp -d -t m-e2e-test.XXX`
trap "rm -rf $TEMPDIR" EXIT HUP INT QUIT TERM

# Build the e2e binary for the tests
echo "--> Building ${CTLNAME} ..."
go build -o $TEMPDIR/m-e2e || exit 1


# Run the tests
echo "--> Running tests"
GOBIN="`which go`"

# If one wants to `go test` a specific package i.e. namespace within 
# this project then `TESTPKG` variable needs to be set with that namespace 
# path before invoking the make command
#
# NOTE:
#   Refer GNUMakefile to see how this script is triggered
if [ -z "${TESTPKG}" ]; 
then
    TESTPKG=$($GOBIN list ./... | grep -v /vendor/)
else
    TESTPKG=$($GOBIN list ./... | grep -v /vendor/ | grep $TESTPKG)
fi

sudo -E PATH=$TEMPDIR:$PATH  -E GOPATH=$GOPATH \
    $GOBIN test ${GOTEST_FLAGS:--cover -timeout=900s} $TESTPKG

