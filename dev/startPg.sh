#!/bin/sh

DATADIR=`mktemp -d /tmp/skedPgData.XXXXXX` || exit 1
docker run -d \
        --name pgDevEnv \
        -e POSTGRES_PASSWORD=pgDevEnv \
        -v ${DATADIR}/:/var/lib/postgresql/data \
        -v ../backend/internal/storage/_dbsetup.sql:/docker-entrypoint-initdb.d/0-dbsetup.sql \
        -v ../backend/internal/storage/testdata/setup.sql:/docker-entrypoint-initdb.d/1-setup.sql\
        -p 54320:5432
        postgres

