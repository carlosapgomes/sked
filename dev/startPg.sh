#echo "${PWD}/data"
docker run -d \
        --name pgDevEnv \
        -e POSTGRES_PASSWORD=pgDevEnv \
        -v ${DATADIR}/:/var/lib/postgresql/data \
        -v ${PWD}/data:/docker-entrypoint-initdb.d/ \
        -p 54320:5432 \
        postgres

        #-v ${STORAGEDIR}/backend/internal/storage/testdata/setup.sql:/docker-entrypoint-initdb.d/1-setup.sql\
