#!/bin/bash

# Start the brindexer service as determined by our environment.

# Environment variables:
# 
#   * Note: Those marked with asterisk are required variables.
#
# For all:
#
#  *POD_TYPE=[cbr]         - Run a br pod .
#  DAEMON_PAUSE=<str>      - If specified and set to "true", then start the
#                            pods but do not start the daemons.  A busy-loop
#                            will be run instead.
#
# For cbr:
#  VERBOSE=<str>          - If specified, enable verbose logging in brindexer.
#

typeset -i errors=0

ARGS=""
DAEMON_PAUSE=true
case "$DAEMON_PAUSE" in
true)
  echo "Do not start the daemon.  Run a busy loop..."
  touch /keep
  while [ -f /keep ]
  do
    sleep 300
  done
  echo "Done."
  exit 0
  ;;
esac


