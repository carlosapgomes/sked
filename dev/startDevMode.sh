#!/bin/sh
# this script sets up the development environment:

# - starts a Postgres docker image with test data
# - compile a backend binary
# - copy the binary to a `/tmp/` folder with the 
#   templates files
# - starts the backend
# - cleans up everything when it receives a `CTRL-C`

# It needs a pre configured basic reverse proxy to map
# a TLS enabled local development DNS entry with the
# following paths:
# - `/api` -> mapping to the backend server port
# - `/` -> mapping to the frontend server port
# See [this blog post](https://carlosapgomes.me/post/localssl/) on
# how to configure a local development with https.



# set the following environment variables before running this script:
# PG_STR
# SENDGRID_API_KEY
# FROM_EMAIL

sigint(){
    echo "signal INT received, script ending";
    docker stop pgDevEnv && docker rm pgDevEnv;
    rm -rf $DATADIR
    rm -rf $SKEDDIR
    exit 0;
}
trap 'sigint' INT TERM
echo
echo "Starting Postgres with a temporary data folder"
echo
DATADIR=`mktemp -d /tmp/skedPgData.XXXXXX` || exit 1
docker run -d \
        --name pgDevEnv \
        -e POSTGRES_PASSWORD=pgDevEnv \
        -v ${DATADIR}/:/var/lib/postgresql/data \
        -v ${PWD}/data:/docker-entrypoint-initdb.d/ \
        -p 54320:5432 \
        postgres || exit 1

echo
echo "Postgres data folder: $DATADIR"
echo
echo "Building the backend binary"
SKEDDIR=`mktemp -d /tmp/skeddir.XXXXXX` || exit 1
cd ../backend
/usr/local/go/bin/go build -o "$SKEDDIR/skedbackend" cmd/main.go;
cd ../dev
echo
echo "Moving assets to temporary folder"
echo
echo "Sked temporary folder: $SKEDDIR"
echo
mkdir "$SKEDDIR/templates" || exit 1;
cp -r ../backend/internal/web/templates/* "$SKEDDIR/templates" || exit 1;
echo "starting the backend"
cd $SKEDDIR
n=0
until [ "$n" -ge 10 ]
do
   echo 
   echo "Waiting postgres to launch..."
   echo
   ./skedbackend && break  
   n=$((n+1)) 
   echo
   echo "Trying again in 5 seconds... (press 'CTRL-C' to abort)"
   echo 
   sleep 5 
done
sigint

