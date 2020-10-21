 -- This is the databse setup script 
 -- execute it with the following command:
 -- psql -U sked -d sked -a -f _dbsetup.sql
CREATE OR REPLACE FUNCTION sked_date_to_char(some_time timestamp with time zone) 
  RETURNS text
AS
$BODY$
    select to_char($1, 'yyyy-mm-dd');
    $BODY$
    LANGUAGE sql
    IMMUTABLE;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  phone TEXT NOT NULL,
  hashedpw TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  active BOOL NOT NULL DEFAULT FALSE,
  email_was_validated BOOL NOT NULL DEFAULT FALSE,
  roles TEXT[]
);
CREATE INDEX users_id_index ON users (id);
CREATE INDEX users_email_index ON users (email);

DROP TABLE IF EXISTS patients;
CREATE TABLE patients (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  name TEXT NOT NULL,
  address TEXT,
  city TEXT,
  state CHAR(2),
  phones TEXT[] NOT NULL,
  created_by UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_by UUID NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX patients_id_index ON patients (id);
CREATE INDEX patients_name_index ON patients (name);

DROP TABLE IF EXISTS appointments;
CREATE TABLE appointments (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  date_time TIMESTAMP WITH TIME ZONE NOT NULL,
  date_txt TEXT, 
  patient_name TEXT NOT NULL,
  patient_id UUID NOT NULL,
  doctor_name TEXT NOT NULL,
  doctor_id UUID NOT NULL,
  notes TEXT,
  canceled BOOL NOT NULL DEFAULT FALSE,
  completed BOOL NOT NULL DEFAULT FALSE,
  created_by UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_by UUID NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX appointments_id_index ON appointments (id);
CREATE INDEX appointments_patient_id_index ON appointments (patient_id);
CREATE INDEX appointments_doctor_id_index ON appointments (doctor_id);
-- CREATE INDEX appointments_date_index ON appointments (date_txt);
CREATE INDEX appointments_date_index ON appointments (sked_date_to_char(date_time));
-- CREATE FUNCTION  insert_date_time_as_text()
 -- RETURNS TRIGGER AS $insert_date_time_as_text$
-- BEGIN
        -- NEW.date_time_txt := date(timezone('UTC', NEW.date_time)); 
        -- RETURN NEW;
-- END;
-- $insert_date_time_as_text$ LANGUAGE PLPGSQL IMMUTABLE;

-- CREATE TRIGGER date_time_text BEFORE INSERT OR UPDATE
    -- ON appointments 
    -- FOR EACH ROW
    -- EXECUTE PROCEDURE insert_date_time_as_text();

DROP TABLE IF EXISTS surgeries;
CREATE TABLE surgeries (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  date_time TIMESTAMP WITH TIME ZONE NOT NULL,
  patient_name TEXT NOT NULL,
  patient_id UUID NOT NULL,
  doctor_name TEXT NOT NULL,
  doctor_id UUID NOT NULL,
  notes TEXT,
  proposed_surgery TEXT NOT NULL,
  canceled BOOL DEFAULT FALSE,
  done BOOL DEFAULT FALSE,
  created_by UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
  updated_by UUID NOT NULL, 
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
CREATE INDEX surgeries_id_index ON surgeries (id);
CREATE INDEX surgeries_patient_id_index ON surgeries (patient_id);
CREATE INDEX surgeries_doctor_id_index ON surgeries (doctor_id);
CREATE INDEX surgeries_date_time_index ON surgeries (sked_date_to_char(date_time));


DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  userid UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE INDEX sessions_id_index ON sessions (id);
CREATE INDEX sessions_userid_index ON sessions (userid);

DROP TABLE IF EXISTS tokens;
CREATE TABLE tokens (
  id TEXT NOT NULL PRIMARY KEY UNIQUE CHECK (id <> ''),
  userid UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  kind TEXT
);
CREATE INDEX tokens_id_index ON tokens (id);
CREATE INDEX tokens_userid_index ON tokens (userid);
