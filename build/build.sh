#!/bin/sh

echo "Building the backend binary"
SKEDDIR=`$PWD/provisioning/roles/backend/files/opt/skedbackend` || exit 1
cd ../backend
/usr/local/go/bin/go build -o "$SKEDDIR/skedbackend" cmd/main.go;
cd ../build
echo "Moving assets to the temporary folder"
mkdir "$SKEDDIR/templates" || exit 1;
cp ../backend/internal/web/templates/* "$SKEDDIR/templates" || exit 1;
echo
echo "Building the frontend"
cd ../frontend/reactfrontend
yarn build
cp -r public/* ../../build/frontend/files/var/www/sked/html/ || exit 1;
cd ../../build

echo
echo "running ansible..."

ansible-playbook -i ansible_inv.yml provisioning/playbook.yml
