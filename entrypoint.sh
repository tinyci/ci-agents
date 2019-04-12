#!/bin/bash

set -e

teardown() {
  service postgresql stop
}

trap teardown INT TERM

[ -f /var/ca/rootCA.pem ] || mkcert --ecdsa --install
[ -f /var/ca/localhost-server.pem ] || mkcert --ecdsa --cert-file /var/ca/localhost-server.pem --key-file /var/ca/localhost-server.key localhost 127.0.0.1 ::1
[ -f /var/ca/localhost-client.pem ] || mkcert --client --ecdsa --cert-file /var/ca/localhost-client.pem --key-file /var/ca/localhost-client.key localhost 127.0.0.1 ::1

if [[ ! -d /var/lib/postgresql/9.6/main ]] || [[ ! -z "${CREATE_DB}" ]]
then
  rm -rf /var/lib/postgresql/9.6/main
  mkdir -p /var/lib/postgresql/9.6/main
  chown postgres:postgres /var/lib/postgresql/9.6/main
  su postgres -c "/usr/lib/postgresql/9.6/bin/initdb -D /var/lib/postgresql/9.6/main"

  service postgresql start

  su postgres -c "bash -c 'createuser -s root'"
  su postgres -c "bash -c 'createuser -s tinyci -P < <(echo -e \"tinyci\ntinyci\n\")'"
  su postgres -c "bash -c 'createdb -O tinyci tinyci'"
else
  service postgresql start
fi

sleep 1
bash -c '/go/bin/migrator -u tinyci -p tinyci migrations/tinyci'

if [ -z "${TESTING}" ]
then
  nginx
fi

"$@"

service postgresql stop
