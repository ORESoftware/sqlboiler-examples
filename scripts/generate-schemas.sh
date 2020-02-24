#!/usr/bin/env bash

set -eo pipefail
cd "$(dirname "$(dirname "$BASH_SOURCE")")"

if ! command -v 'schemats'; then
    npm install -g schemats
fi


schemats generate -c 'postgres://postgres:postgres@localhost/postgres' -t 'user_table' -o schemas/user.ts

schemats generate -c 'postgres://postgres:postgres@localhost/postgres' -t 'user_map_table' -o schemas/user-map.ts
