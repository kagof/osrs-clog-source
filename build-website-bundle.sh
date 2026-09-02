#!/bin/bash
set -o pipefail

BASE_DIR=$(dirname "$BASH_SOURCE")
echo "using base dir $BASE_DIR"

if [ -z "$BASE_DIR" ]; then
    echo "base dir is empty string"
    exit 1
fi

rm -rf "$BASE_DIR/dist"

mkdir -p "$BASE_DIR/dist"

cp -r "$BASE_DIR/website/assets" "$BASE_DIR/dist/assets"
cp "$BASE_DIR/website/index.html" "$BASE_DIR/dist/index.html"
cp "$BASE_DIR/LICENSE" "$BASE_DIR/dist/LICENSE.txt"