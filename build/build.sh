#!/bin/sh

echo "Building the backend binary"
SKEDDIR=`$PWD/provisioning` || exit 1
cd ../backend
/usr/local/go/bin/go build -o "$SKEDDIR/skedbackend" cmd/main.go;
cd ../dev
echo "Moving assets to the temporary folder"
echo "Sked temporary folder: $SKEDDIR"
mkdir "$SKEDDIR/templates" || exit 1;
cp ../backend/internal/web/templates/* "$SKEDDIR/templates" || exit 1;
echo "starting the backend"
echo "Waiting postgres to launch..."
cd $SKEDDIR
n=0
until [ "$n" -ge 10 ]
do
   ./skedbackend && break  
   n=$((n+1)) 
   sleep 5 
done

