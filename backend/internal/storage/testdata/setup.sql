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


CREATE TABLE patients (
  id UUID NOT NULL PRIMARY KEY UNIQUE,
  name VARCHAR(100),
  address VARCHAR(150),
  city VARCHAR(100),
  state CHAR(2),
  phones TEXT[],
  created_by UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_by UUID NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE appointments{
    id UUID NOT NULL PRIMARY KEY UNIQUE,
    date_time TIMESTAMPWITH TIME ZONE,
    patient_name VARCHAR(100),
    patient_id UUID NOT NULL, 
    doctor_name VARCHAR(100),
    doctor_id UUID NOT NULL,
    notes TEXT,
    canceled BOOL,
    completed BOOL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
};

CREATE TABLE surgeries{
    id UUID NOT NULL PRIMARY KEY UNIQUE,
    date_time TIMESTAMPWITH TIME ZONE,
    patient_name VARCHAR(100),
    patient_id UUID NOT NULL, 
    doctor_name VARCHAR(100),
    doctor_id UUID NOT NULL,
    notes TEXT,
    proposed_surgery TEXT,
    canceled BOOL,
    done BOOL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_by UUID NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE
};

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

\i insert_tokes.sql
\i insert_sessions.sql
\i insert_users.sql
\i insert_appointments.sql
\i insert_surgeries.sql

