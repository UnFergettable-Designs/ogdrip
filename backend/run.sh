#!/bin/bash
# Activate Flox environment
eval "$(flox activate)"

# Run the server with the corrected URL
go run server.go -url="https://valkyriefitnessvn.com" 