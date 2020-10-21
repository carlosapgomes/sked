DATADIR=`mktemp -d /tmp/skedPgData.XXXXXX` || exit 1
docker run -d \
        --name pgDevEnv \
        -e POSTGRES_PASSWORD=pgDevEnv \
        -v ${DATADIR}/:/var/lib/postgresql/data \
        -v ${PWD}/data:/docker-entrypoint-initdb.d/ \
        -p 54320:5432 \
        postgres

