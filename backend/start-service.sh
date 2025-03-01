#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
  echo "Loaded environment variables:"
  echo "PORT=$PORT"
  echo "ENABLE_CORS=$ENABLE_CORS"
  echo "BASE_URL=$BASE_URL"
  echo "OUTPUT_DIR=$OUTPUT_DIR"
fi

# Start the backend service
go run main.go server.go service.go -service 