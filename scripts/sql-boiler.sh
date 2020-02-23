#!/usr/bin/env bash

set -eo pipefail
cd "$(dirname "$(dirname "$BASH_SOURCE")")"

if [[ -d 'github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql' ]]; then
    go get 'github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql'
    go install "$GOPATH/src/github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql"
fi

export PATH="${GOPATH}/bin:$PATH"


export db_name="${sqlboiler_db_name:-postgres}"

if [[ -z "$db_pwd" ]]; then
    export db_pwd="${sqlboiler_db_pwd}"
fi


if [[ -z "$db_pwd" ]]; then
  echo 'db_pwd is empty - please use an env variable either "$db_pwd" or "$sqlboiler_db_pwd" ' > /dev/stderr
  exit 1;
fi


export db_user="${sqlboiler_db_user:-postgres}"
export db_host="${sqlboiler_db_host:-localhost}"


node sqlboiler.js > sqlboiler.js.json

sqlboiler -c 'sqlboiler.js.json' psql