#! /bin/bash

set -eu

if [ $# != 0 ]; then
    echo "Usage: $0"
    exit 1
fi

BASEDIR=$(dirname $0)

# Sign the auth credential
okp4d --keyring-backend test --keyring-dir ${BASEDIR}/../example credential sign --from exec-svc "${BASEDIR}/../example/vc-exec-auth.jsonld"
