#!/bin/sh

# Use this script to run your program LOCALLY.

set -e # Exit early if any commands fail

(
  cd "$(dirname "$0")" # Ensure compile steps are run within the repository directory
  go build -o /tmp/build-shell-go app/*.go
)

# Copied from .codecrafters/run.sh
#
# - Edit this to change how your program runs locally
# - Edit .codecrafters/run.sh to change how your program runs remotely
exec /tmp/build-shell-go "$@"
