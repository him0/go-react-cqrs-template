#!/bin/bash

set -e

echo "Generating Go server code from OpenAPI spec..."

# Check if openapi-generator-cli is installed
if ! command -v openapi-generator-cli &> /dev/null
then
    echo "openapi-generator-cli is not installed."
    echo "Please install it with: npm install -g @openapitools/openapi-generator-cli"
    exit 1
fi

# Generate Go server code
openapi-generator-cli generate \
  -i openapi/openapi.yaml \
  -g go-server \
  -o pkg/generated \
  --additional-properties=packageName=openapi,router=chi,sourceFolder=.

echo "API code generation completed!"
