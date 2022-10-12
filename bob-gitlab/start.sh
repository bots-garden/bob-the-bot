#!/bin/bash

BOT_NAME="@swannou" \
BOT_TOKEN="${SWANNOU_TOKEN}" \
API_URL="https://gitlab.com/api/v4" \
capsule -wasm=./bob.wasm -mode=http -httpPort=8080

#https://gitlab.com/k33g_org/bob-the-bot
# REPO_NAME="bob-the-bot" \
# REPO_OWNER="k33g_org" \