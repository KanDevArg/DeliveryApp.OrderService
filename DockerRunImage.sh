#!/bin/bash
echo "${BASH_VERSION}"
FOLDER_NAME=$(basename $(pwd))
CONTAINER_NAME=$(echo "$FOLDER_NAME" |  tr '[:upper:]' '[:lower:]' )


echo "Running image..."$CONTAINER_NAME
docker run -p 8080:8080 $CONTAINER_NAME