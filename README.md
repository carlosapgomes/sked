# Sked(duler)

This program will keep track of patients appointments for a specific
specialty/clinic.

## Use Cases

- Only an admin can add a new user
- A user can schedule a patient for any available date
- A user can insert a new patient
- A user can modify a patient's phoneNumber
- Only an admin can modify a patient's name
- A user can mark an appointment as removed only until the scheduled date
- A user can unmark an appointment as removed only until the scheduled date
- A user can mark an appointment as complete only at its scheduled date
- A user can unmark an appointment as complete only at its scheduled date
- Nobody can change an appointment's date. A new appointment should be created.

## Entities

### User

- id (uuidV4)
- Name
- PhoneNumber
- isAdmin?
- createdBy
- updatedBy
- createdAt

### Patients

- id (uuidV4)
- Name
- PhoneNumber
- createdBy
- updatedBy
- createdAt

### Appointments

- id (uuidV4)
- Date/Time
- Doctor
- Patient
- Notes
- isCancelled?
- completed?
- createdBy
- updatedBy
- createdAt

## Backend

### Golang Http server

#### Use server side sessions

### Postgres

Tables:

- Users
- Patients
- Appointments
- ActivitiesLog

## FrontEnd

### Client side rendering

#### Lit-element pages:

- home/login page
- add/edit appointment
- add/edit patient
- add/edit user
- view scheduling
- view activities log
