#! /bin/bash

set -eu

if [ $# != 0 ]; then
    echo "Usage: $0"
    exit 1
fi

BASEDIR=$(dirname "$0")

# Sign the auth credential
axoned --keyring-backend test --keyring-dir "${BASEDIR}"/../example credential sign --from exec-svc --purpose authentication "${BASEDIR}/../example/vc-exec-auth.jsonld"
