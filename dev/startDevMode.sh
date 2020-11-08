#!/bin/sh
# this script sets up the development environment:

# - starts a Postgres docker image with test data
# - compile a backend binary
# - copy the binary to a `/tmp/` folder with the 
#   templates files
# - starts the backend
# - runs the frontend server
# - cleans up everything when it receives a `CTRL-C`

# It needs a pre configured basic reverse proxy to map
# a TLS enabled local development dns entry with the
# following paths:
# - `/api` -> mapping to the backend server port
# - `/` -> mapping to the frontend server port

# set the following environment variables:
# PG_STR
# SENDGRID_API_KEY
# FROM_EMAIL
echo "Setting environment variables"
source ./setEnv.sh

sigint(){
    echo "signal INT received, script ending";
    docker stop pgDevEnv && docker rm pgDevEnv;
    rm -rf $DATADIR
    rm -rf $SKEDDIR
    pkill -f yarn
    exit 0;
}
trap sigint SIGINT
CURRENTDIR=$PWD
echo "Starting Postgres with a temporary data folder"
DATADIR=`mktemp -d /tmp/skedPgData.XXXXXX` || exit 1
docker run -d \
        --name pgDevEnv \
        -e POSTGRES_PASSWORD=pgDevEnv \
        -v ${DATADIR}/:/var/lib/postgresql/data \
        -v ${PWD}/data:/docker-entrypoint-initdb.d/ \
        -p 54320:5432 \
        postgres || exit 1

echo "Postgres data folder: $DATADIR"
echo "Building the backend binary"
SKEDDIR=`mktemp -d /tmp/skeddir.XXXXXX` || exit 1
cd ../backend
/usr/local/go/bin/go build -o "$SKEDDIR/skedbackend" cmd/main.go;
cd ../dev
echo "Moving assets to the temporary folder"
echo "Sked temporary folder: $SKEDDIR"
mkdir "$SKEDDIR/templates" || exit 1;
cp ../backend/internal/web/templates/* "$SKEDDIR/templates" || exit 1;
echo "starting the frontend"
cd $CURRENTDIR
cd ../frontend/reactfrontend
yarn start  &
echo "Waiting postgres to launch..."
cd $SKEDDIR
n=0
until [ "$n" -ge 8 ]
do
   ./skedbackend && break  # substitute your command here
   n=$((n+1)) 
   sleep 5 
done


