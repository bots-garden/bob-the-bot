#!/bin/bash

BOT_NAME="@jayisabot" \
BOT_TOKEN="${JAYBOT_TOKEN}" \
REPO_NAME="bob-the-bot" \
REPO_OWNER="bots-garden" \
API_URL="https://api.github.com" \
capsule -wasm=./bob.wasm -mode=http -httpPort=8080
