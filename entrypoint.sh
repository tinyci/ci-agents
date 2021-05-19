#!/bin/bash

set -euo pipefail

POSTGRES_VERSION=12

cat welcome.ans

teardown() {
  service postgresql stop
}

if [ "$(id -u)" == "0" ]
then
  trap teardown INT TERM

  [ -f /var/ca/rootCA.pem ] || mkcert --ecdsa --install
  [ -f /var/ca/localhost-server.pem ] || mkcert --ecdsa --cert-file /var/ca/localhost-server.pem --key-file /var/ca/localhost-server.key localhost 127.0.0.1 ::1
  [ -f /var/ca/localhost-client.pem ] || mkcert --client --ecdsa --cert-file /var/ca/localhost-client.pem --key-file /var/ca/localhost-client.key localhost 127.0.0.1 ::1

  config=/etc/postgresql/${POSTGRES_VERSION}/main/postgresql.conf

  grep -iv shared_buffers ${config} >${config}.tmp && mv ${config}.tmp ${config}
  echo shared_buffers=2GB >> ${config}

  if [[ ! -d /var/lib/postgresql/${POSTGRES_VERSION}/main ]] || [[ ! -z "${CREATE_DB}" ]]
  then
    rm -rf /var/lib/postgresql/${POSTGRES_VERSION}/main
    mkdir -p /var/lib/postgresql/${POSTGRES_VERSION}/main
    chown postgres:postgres /var/lib/postgresql/${POSTGRES_VERSION}/main
    su postgres -c "/usr/lib/postgresql/${POSTGRES_VERSION}/bin/initdb -D /var/lib/postgresql/${POSTGRES_VERSION}/main"

    service postgresql start

    su postgres -c "bash -c 'createuser -s root'"
    su postgres -c "bash -c 'createuser -s tinyci -P < <(echo -e \"tinyci\ntinyci\n\")'"
    su postgres -c "bash -c 'createdb -O tinyci tinyci'"
  else
    service postgresql start
  fi

  if [ -z "${TESTING}" ]
  then
    if go install -v ./... 
    then
      tinyci migrate
    else
      echo >&2 "Code didn't compile, not migrating!"
    fi

    caddy start -config /Caddyfile -watch
  else
    if go install -v ./... 
    then
      tinyci --config .config/services.yaml.example migrate
    else
      echo >&2 "Code didn't compile, not migrating!"
    fi
  fi
fi

exec "$@"
