#! /bin/bash

set -eu

if [ $# != 2 ]; then
    echo "Usage: $0 [sender_addr] [dataverse_addr]"
    exit 1
fi

BASEDIR=$(dirname $0)

./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-exec-order.jsonld" user $1 $2
./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-exec.jsonld" exec-svc $1 $2
