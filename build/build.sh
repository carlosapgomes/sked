#!/bin/sh

echo "Preparing folders"
if [ ! -d  "$PWD/provisioning/roles/backend/files/opt/skedbackend" ]; then
    echo "making folders"
    mkdir -p "$PWD/provisioning/roles/backend/files/opt/skedbackend" || exit 1
fi
SKEDDIR="$PWD/provisioning/roles/backend/files/opt/skedbackend" || exit 1
echo "changing dir to ../backend"
cd ../backend
echo "Building skedbackend"
/usr/local/go/bin/go build -o "$SKEDDIR/skedbackend" cmd/main.go
cd ../build
echo "Moving assets to the temporary folder"
if [ ! -d  "$SKEDDIR/templates" ]; then
    mkdir "$SKEDDIR/templates" || exit 1
fi
cp -r ../backend/internal/web/templates/* "$SKEDDIR/templates" || exit 1
echo
echo "Building the frontend"
cd ../frontend/reactfrontend
yarn build
echo "cd to project root"
cd ../../
if [ ! -d "$PWD/build/provisioning/roles/frontend/files/var/www/sked/html" ]; then
    echo "making frontend folders"
    mkdir -p "$PWD/build/provisioning/roles/frontend/files/var/www/sked/html" || exit 1
fi
echo "current folder: $PWD"
cp -R "$PWD/frontend/reactfrontend/build/." "$PWD/build/provisioning/roles/frontend/files/var/www/sked/html" || exit 1
cd build
echo
echo "running ansible..."

ansible-playbook -i ansible_inv.yml provisioning/playbook.yml
