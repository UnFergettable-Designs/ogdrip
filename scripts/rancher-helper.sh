#!/bin/bash

# Helper script for Rancher Desktop operations with nerdctl
# Usage: ./rancher-helper.sh [command]

set -e

check_status() {
  if pgrep -f "Rancher Desktop" > /dev/null; then
    echo "Rancher Desktop is running"
    return 0
  else
    echo "Rancher Desktop is not running"
    return 1
  fi
}

check_nerdctl() {
  if command -v nerdctl &> /dev/null; then
    echo "nerdctl is installed"
    return 0
  else
    echo "nerdctl is not installed. Run 'brew install nerdctl' to install it."
    return 1
  fi
}

case "$1" in
  "start")
    echo "Starting Rancher Desktop..."
    open -a "Rancher Desktop"
    ;;

  "stop")
    echo "Stopping Rancher Desktop..."
    pkill -f "Rancher Desktop" || echo "Rancher Desktop not running"
    ;;

  "status")
    check_status
    ;;

  "ensure-running")
    if ! check_status; then
      echo "Starting Rancher Desktop..."
      open -a "Rancher Desktop"
      echo "Waiting for Rancher Desktop to be ready..."
      sleep 10  # Give it some time to start
    fi
    ;;

  "build")
    shift
    ./scripts/rancher-helper.sh ensure-running
    check_nerdctl || exit 1
    echo "Building with nerdctl: $@"
    nerdctl build $@
    ;;

  "run")
    shift
    ./scripts/rancher-helper.sh ensure-running
    check_nerdctl || exit 1
    echo "Running with nerdctl: $@"
    nerdctl run $@
    ;;

  "check-nerdctl")
    check_nerdctl
    if [ $? -eq 0 ]; then
      echo "Testing nerdctl connection to Rancher Desktop..."
      nerdctl info &> /dev/null
      if [ $? -eq 0 ]; then
        echo "nerdctl is properly configured with Rancher Desktop"
      else
        echo "nerdctl may not be properly configured with Rancher Desktop. Make sure Rancher Desktop is running and using containerd runtime."
      fi
    fi
    ;;

  *)
    echo "Unknown command: $1"
    echo "Available commands: start, stop, status, ensure-running, build, run, check-nerdctl"
    exit 1
    ;;
esac
