#! /bin/bash

set -eu

if [ $# != 2 ]; then
    echo "Usage: $0 [sender_addr] [dataverse_addr]"
    exit 1
fi

BASEDIR=$(dirname $0)

./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-s3-desc.jsonld" s3-issuer $1 $2
sleep 10
./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-s3-gov.jsonld" s3-issuer $1 $2
sleep 10
./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-data-desc.jsonld" data-issuer $1 $2
sleep 10
./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-data-gov.jsonld" data-issuer $1 $2
sleep 10
./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-publish.jsonld" s3-issuer $1 $2
sleep 10
./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-exec-order.jsonld" user $1 $2
sleep 10
./${BASEDIR}/submit-vc.sh "${BASEDIR}/../example/vc-exec.jsonld" exec-svc $1 $2
