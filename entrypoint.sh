#!/bin/bash

set -e

POSTGRES_VERSION=11

cat welcome.ans

teardown() {
  service postgresql stop
}

trap teardown INT TERM

[ -f /var/ca/rootCA.pem ] || mkcert --ecdsa --install
[ -f /var/ca/localhost-server.pem ] || mkcert --ecdsa --cert-file /var/ca/localhost-server.pem --key-file /var/ca/localhost-server.key localhost 127.0.0.1 ::1
[ -f /var/ca/localhost-client.pem ] || mkcert --client --ecdsa --cert-file /var/ca/localhost-client.pem --key-file /var/ca/localhost-client.key localhost 127.0.0.1 ::1

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

while ! bash -c '/go/bin/migrator -q -u tinyci -p tinyci migrations/tinyci'
do 
  sleep 1
  i=$(($i + 1));
  if [ "$i" -gt 10 ]
  then
    echo >&2 Timed out
    exit 1
  fi
done

if [ -z "${TESTING}" ]
then
  nginx
fi

"$@"

service postgresql stop
