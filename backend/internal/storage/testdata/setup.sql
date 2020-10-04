-- CREATE TABLE users (
  -- id UUID NOT NULL PRIMARY KEY UNIQUE,
  -- name TEXT,
  -- email TEXT NOT NULL UNIQUE,
  -- phone TEXT,
  -- hashedpw TEXT NOT NULL,
  -- created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  -- updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  -- active BOOL NOT NULL DEFAULT FALSE,
  -- email_was_validated BOOL NOT NULL DEFAULT FALSE,
  -- roles TEXT[]
-- );
-- CREATE INDEX users_id_index ON users (id);
-- CREATE INDEX users_email_index ON users (email);


-- CREATE TABLE patients (
  -- id UUID NOT NULL PRIMARY KEY UNIQUE,
  -- name TEXT,
  -- address TEXT, 
  -- city TEXT, 
  -- state CHAR(2),
  -- phones TEXT[],
  -- created_by UUID NOT NULL,
  -- created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  -- updated_by UUID NOT NULL,
  -- updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
-- );
-- CREATE INDEX patients_id_index ON patients (id);
-- CREATE INDEX patients_name_index ON patients (name);

-- CREATE TABLE appointments{
    -- id UUID NOT NULL PRIMARY KEY UNIQUE,
    -- date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    -- patient_name TEXT NOT NULL, 
    -- patient_id UUID NOT NULL, 
    -- doctor_name TEXT NOT NULL, 
    -- doctor_id UUID NOT NULL,
    -- notes TEXT,
    -- canceled BOOL NOT NULL DEFAULT FALSE,
    -- completed BOOL NOT NULL DEFAULT FALSE,
    -- created_by UUID NOT NULL,
    -- created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- updated_by UUID NOT NULL,
    -- updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
-- };
-- CREATE INDEX appointments_id_index ON appointments (id);
-- CREATE INDEX appointments_patient_id_index ON appointments (patient_id);
-- CREATE INDEX appointments_doctor_id_index ON appointments (doctor_id);
-- CREATE INDEX appointments_date_index ON appointments (date(date_time));

-- CREATE TABLE surgeries{
    -- id UUID NOT NULL PRIMARY KEY UNIQUE,
    -- date_time TIMESTAMP WITH TIME ZONE NOT NULL,
    -- patient_name TEXT NOT NULL, 
    -- patient_id UUID NOT NULL, 
    -- doctor_name TEXT NOT NULL, 
    -- doctor_id UUID NOT NULL,
    -- notes TEXT,
    -- proposed_surgery TEXT NOT NULL,
    -- canceled BOOL DEFAULT FALSE,
    -- done BOOL DEFAULT FALSE,
    -- created_by UUID NOT NULL,
    -- created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- updated_by UUID NOT NULL,
    -- updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
-- };
-- CREATE INDEX surgeries_id_index ON surgeries (id);
-- CREATE INDEX surgeries_patient_id_index ON surgeries (patient_id);
-- CREATE INDEX surgeries_doctor_id_index ON surgeries (doctor_id);
-- CREATE INDEX surgeries_date_index ON surgeries (date(date_time));

-- CREATE TABLE sessions (
  -- id UUID NOT NULL PRIMARY KEY UNIQUE,
  -- userid UUID NOT NULL,
  -- created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  -- expires_at TIMESTAMP WITH TIME ZONE NOT NULL
-- );
-- CREATE INDEX sessions_id_index ON sessions (id);
-- CREATE INDEX sessions_userid_index ON sessions (userid);

-- CREATE TABLE tokens (
  -- id TEXT NOT NULL PRIMARY KEY UNIQUE CHECK (id <> ''),
  -- userid UUID NOT NULL,
  -- created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  -- expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  -- kind TEXT
-- );
-- CREATE INDEX tokens_id_index ON tokens (id);
-- CREATE INDEX tokens_userid_index ON tokens (userid);
\i insert_tokens.sql
\i insert_sessions.sql
\i insert_users.sql
\i insert_patients.sql
\i insert_appointments.sql
\i insert_surgeries.sql

