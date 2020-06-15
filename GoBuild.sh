#!/bin/bash
PRE_PWD=$(pwd)
FOLDER_NAME=$(basename $(pwd))

echo "Building main.go ..."

mkdir -p bin

go build -o $PRE_PWD/bin/$FOLDER_NAME .

echo "Done."
