#!/bin/bash
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "sked" <<'EOF'
INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at,updated_at,active, email_was_validated, roles) VALUES(
'dcce1beb-aee6-4a4d-b724-94d470817323',
'Alice Jones',
'alice@example.com','6544334535',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-23 17:25:00 UTC',
'2019-06-23 17:25:00 UTC',
TRUE,FALSE,
'{"Clerk"}'
);
INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at, updated_at, active, email_was_validated, roles) VALUES(
		'68b1d5e2-39dd-4713-8631-a08100383a0f','Bob',
    'bob@example.com','5434534534',
    '$2a$12$kHna5vstSSusP9VFC89tZ.317kInW7dZRL8snvnAej66wgQnyaQte',
    '2019-06-24 17:25:00 UTC',
    '2019-06-24 17:25:00 UTC',
		TRUE,TRUE,
		'{"Clerk","Admin"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at, updated_at, active, email_was_validated, roles) VALUES(
 '85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd','Valid User','valid@user.com',
 '6544332135','$2a$12$kHna5vstSSusP9VFC89tZ.317kInW7dZRL8snvnAej66wgQnyaQte',
 '2019-06-25 17:25:00 UTC',
 '2019-06-25 17:25:00 UTC',
 TRUE,TRUE,
 '{"Clerk","Admin"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at, updated_at, active, email_was_validated, roles) VALUES(
'ecadbb28-14e6-4560-8574-809c6c54b9cb','Barack Obama','bobama@somewhere.com',
'6544332135',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-26 17:25:00 UTC',
'2019-06-26 17:25:00 UTC',
FALSE,TRUE,
'{"Clerk"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at, updated_at, active, email_was_validated, roles) VALUES(
'ca16fc9d-df7b-4594-97e3-264432145b01','SpongeBob Squarepants','spongebob@somewhere.com',
'65949340','$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-27 17:25:00 UTC',
'2019-06-27 17:25:00 UTC',
FALSE,TRUE,
'{"Clerk"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at, updated_at, active, email_was_validated, roles) VALUES(
'27f9802b-acb3-4852-bf97-c4ed4c3b3658',
'Tim Berners-Lee',
'tblee@somewhere.com',
'0323949324',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-28 17:25:00 UTC',
'2019-06-28 17:25:00 UTC',
FALSE,TRUE,
'{"Clerk"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at, updated_at, active, email_was_validated, roles) VALUES(
'f06244b9-97e5-4f1a-bae0-3b6da7a0b604',
'Dr. House',
'house@doctor.com',
'6544332135',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-28 17:25:00 UTC',
'2019-06-28 17:25:00 UTC',
TRUE,TRUE,
'{"Doctor"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, 
    created_at, updated_at, active, email_was_validated, roles) VALUES(
'a520df95-02fa-4d86-8eef-58385c354b29',
'Shaun Murphy',
'shaun@thegooddoctor.com',
'64532332135',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-28 17:25:00 UTC',
'2019-06-28 17:25:00 UTC',
TRUE,TRUE,
'{"Doctor"}'
);

EOF

