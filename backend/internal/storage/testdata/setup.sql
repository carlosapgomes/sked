CREATE TABLE users (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  name TEXT,
  email TEXT NOT NULL UNIQUE,
  phone TEXT,
  hashedpw TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  active BOOL NOT NULL DEFAULT FALSE,
  email_was_validated BOOL NOT NULL DEFAULT FALSE,
  roles TEXT[]
);


-- CREATE TABLE patients (
--   id UUID NOT NULL PRIMARY KEY UNIQUE,
--   name VARCHAR(100),
--   email VARCHAR(254) NOT NULL UNIQUE,
--   address VARCHAR(150),
--   city VARCHAR(100),
--   state CHAR(2),
--   phones TEXT[],
--   created_at TIMESTAMP WITH TIME ZONE,
--   updated_at TIMESTAMP WITH TIME ZONE
-- );

CREATE TABLE sessions (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  userid UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  expires_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE tokens (
  id TEXT NOT NULL PRIMARY KEY UNIQUE CHECK (id <> ''),
  userid UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  expires_at TIMESTAMP WITH TIME ZONE,
  kind TEXT
);


-- Alice's password: 'test1234'
INSERT INTO users (id, name, email, phone, hashedpw, created_at,updated_at,active, email_was_validated, roles) VALUES(
'dcce1beb-aee6-4a4d-b724-94d470817323',
'Alice Jones',
'alice@example.com','6544334535',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-23 17:25:00 UTC',
'2019-06-23 17:25:00 UTC',
TRUE,FALSE,
'{"Common"}'
);

-- Bob'b password: 'test1234'
INSERT INTO users (id, name, email, phone, hashedpw, created_at, updated_at, active, email_was_validated, roles) VALUES(
		'68b1d5e2-39dd-4713-8631-a08100383a0f','Bob',
    'bob@example.com','6544334535',
    '$2a$12$kHna5vstSSusP9VFC89tZ.317kInW7dZRL8snvnAej66wgQnyaQte',
    '2019-06-24 17:25:00 UTC',
    '2019-06-24 17:25:00 UTC',
		TRUE,TRUE,
		'{"Common","Admin"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, created_at, updated_at, active, email_was_validated, roles) VALUES(
 '85f45ff9-d31c-4ff7-94ac-5afb5a1f0fcd','Valid User','valid@user.com',
 '6544332135','$2a$12$kHna5vstSSusP9VFC89tZ.317kInW7dZRL8snvnAej66wgQnyaQte',
 '2019-06-25 17:25:00 UTC',
 '2019-06-25 17:25:00 UTC',
 TRUE,TRUE,
 '{"Common","Admin"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, created_at, updated_at, active, email_was_validated, roles) VALUES(
'ecadbb28-14e6-4560-8574-809c6c54b9cb','Barack Obama','bobama@somewhere.com',
'6544332135',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-26 17:25:00 UTC',
'2019-06-26 17:25:00 UTC',
FALSE,TRUE,
'{"Common"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, created_at, updated_at, active, email_was_validated, roles) VALUES(
'ca16fc9d-df7b-4594-97e3-264432145b01','SpongeBob Squarepants','spongebob@somewhere.com',
'65949340','$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-27 17:25:00 UTC',
'2019-06-27 17:25:00 UTC',
FALSE,TRUE,
'{"Common"}'
);

INSERT INTO users (id, name, email, phone, hashedpw, created_at, updated_at, active, email_was_validated, roles) VALUES(
'27f9802b-acb3-4852-bf97-c4ed4c3b3658',
'Tim Berners-Lee',
'tblee@somewhere.com',
'0323949324',
'$2a$12$I9BW22CbzLHzY9ORTRhkEuEtq8ufJVMf1dX9CKFlo4W9cIaAjD0Je',
'2019-06-28 17:25:00 UTC',
'2019-06-28 17:25:00 UTC',
FALSE,TRUE,
'{"Common"}'
);

UPDATE users SET phone='5434534534' where id='68b1d5e2-39dd-4713-8631-a08100383a0f';


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