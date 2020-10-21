#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "sked" <<'EOF'
INSERT INTO sessions (id, userid, created_at, expires_at) VALUES(
'144f1223-70f4-4fca-9c99-c58eed7a0f4a',
'dcce1beb-aee6-4a4d-b724-94d470817323',
'2019-06-23 17:00:00 UTC',
'2019-06-23 17:00:00 UTC'
);

INSERT INTO sessions (id, userid, created_at, expires_at) VALUES(
'ddfd72b2-d57b-4fc4-9265-55b9dae7287c',
'dcce1beb-aee6-4a4d-b724-94d470817323',
'2019-06-23 18:25:22 UTC',
'2019-06-23 18:45:22 UTC'
);
EOF
