#!/bin/bash

# Set env variables defined in .env
ENV_FILE=".env"

# Read the .env file and set the variables
while IFS='=' read -r key value; do
  export "$key=$value"
done < "$ENV_FILE"

# Build and run executable
cd front && go build && ./my-wolverine-events-front