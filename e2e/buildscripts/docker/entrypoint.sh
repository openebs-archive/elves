#!/bin/sh

set -ex

# Launch m-e2e
exec /usr/local/bin/m-e2e run 1>&2
