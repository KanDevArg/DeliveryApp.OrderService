#!/bin/bash
echo "${BASH_VERSION}"
FOLDER_NAME=$(basename $(pwd))
CONTAINER_NAME=$(echo "$FOLDER_NAME" |  tr '[:upper:]' '[:lower:]' )

echo "Docker image prune"
docker image prune

echo "Remove image"
docker rmi -f Image $CONTAINER_NAME

echo "Building container "$CONTAINER_NAME 
docker build -t $CONTAINER_NAME .

