/* This is the setup script for Gobackend Database
 run this script with the following command:
 psql -U sked -d sked -a -f _dbsetup.sql
*/
DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  name TEXT,
  email TEXT NOT NULL UNIQUE,
  phone TEXT,
  hashedpw CHAR(60) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  active BOOL NOT NULL DEFAULT FALSE,
  email_was_validated BOOL NOT NULL DEFAULT FALSE,
  roles TEXT[]
);
CREATE INDEX users_id_index ON users (id);
CREATE INDEX users_email_index ON users (email);

DROP TABLE IF EXISTS patients;
CREATE TABLE patients (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  name TEXT,
  address TEXT,
  city TEXT,
  state CHAR(2),
  phones TEXT[],
  created_by UUID,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_by UUID,
  updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX patients_id_index ON patients (id);
CREATE INDEX patients_name_index ON patients (name);

DROP TABLE IF EXISTS appointments;
CREATE TABLE appointments (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  date_time TIMESTAMP WITH TIME ZONE,
  patient_name TEXT,
  patient_id UUID,
  doctor_name TEXT,
  doctor_id UUID,
  notes TEXT,
  canceled BOOL,
  completed BOOL,
  created_by UUID,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_by UUID,
  updated_at TIMESTAMP WITH TIME ZONE
)
CREATE INDEX appointments_id_index ON appointments (id);
CREATE INDEX appointments_patient_id_index ON appointments (patient_id);
CREATE INDEX appointments_doctor_id_index ON appointments (doctor_id);
CREATE INDEX appointments_date_time_index ON appointments (date_time);

DROP TABLE IF EXISTS surgeries;
CREATE TABLE surgeries (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  date_time TIMESTAMP WITH TIME ZONE,
  patient_name TEXT,
  patient_id UUID,
  doctor_name TEXT,
  doctor_id UUID,
  notes TEXT,
  proposed_surgery TEXT,
  canceled BOOL,
  done BOOL,
  created_by UUID,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_by UUID,
  updated_at TIMESTAMP WITH TIME ZONE
)
CREATE INDEX surgeries_id_index ON surgeries (id);
CREATE INDEX surgeries_patient_id_index ON surgeries (patient_id);
CREATE INDEX surgeries_doctor_id_index ON surgeries (doctor_id);
CREATE INDEX surgeries_date_time_index ON surgeries (date_time);


DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  userid UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  expires_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sessions_id_index ON sessions (id);
CREATE INDEX sessions_userid_index ON sessions (userid);

DROP TABLE IF EXISTS tokens;
CREATE TABLE tokens (
  id TEXT NOT NULL PRIMARY KEY UNIQUE CHECK (id <> ''),
  userid UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  expires_at TIMESTAMP WITH TIME ZONE,
  kind TEXT
);

CREATE INDEX tokens_id_index ON tokens (id);
CREATE INDEX tokens_userid_index ON tokens (userid);
