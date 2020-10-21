#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "sked" <<'EOF'
INSERT INTO tokens (id, userid, created_at, expires_at, kind) VALUES(
  '7telsDIFzlzIlZ1fiH8pDtoFiJMoBUi69j9525jt',
  'dcce1beb-aee6-4a4d-b724-94d470817323',
  '2019-06-23 18:25:00 UTC',
  '2019-06-23 18:55:00 UTC',
  'PwReset'
);

INSERT INTO tokens (id, userid, created_at, expires_at, kind) VALUES(
  '7FXsSqU5UC6I9632BndMkSFDCDNO4i1Z83v9KGd2',
  'dcce1beb-aee6-4a4d-b724-94d470817323',
  '2019-06-23 18:25:00 UTC',
  '2019-06-23 19:25:00 UTC',
  'ValidateEmail'
);
EOF
