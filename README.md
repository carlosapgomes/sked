# Sked(duler)

This program will keep track of patients appointments/surgeries for a specific
specialty/clinic.

## Use Cases

- Only an admin can add a new user
- A user with role 'doctor' can only schedule his/her own appointments or surgeries
- A user with role 'common' can schedule appointments and surgeries for doctors
- A user with role 'common' does not have schedules
- A user can schedule a patient for any available date
- A user can schedule a surgery
- A user can insert a new patient
- A user can modify a patient's phoneNumber
- Only an admin can modify a patient's name
- A user can mark an appointment/surgery as removed only until the scheduled date
- A user can unmark an appointment/surgery as removed only until the scheduled date
- A user can mark an appointment/surgery as complete only at its scheduled date
- A user can unmark an appointment/surgery as complete only at its scheduled date
- Nobody can change an appointment/surgery's date. A new appointment/surgery
should be created.

## Entities

### User

- id (uuidV4)
- Name
- Phone
- isAdmin?
- createdBy
- createdAt
- updatedBy
- updatedAt

### Patients

- id (uuidV4)
- Name
- Address
- City
- State
- Phones
- createdBy
- createdAt
- updatedBy
- updatedAt

### Appointments

- id (uuidV4)
- Date/Time
- Doctor (uuidV4)
- Patient (uuidV4)
- Notes
- isCancelled?
- completed?
- createdBy
- createdAt
- updatedBy
- updatedAt

### Surgeries

- id (uuidV4)
- Date/Time
- Patient (uuidV4)
- Doctor (uuidV4)
- Notes
- Proposed surgery
- Cancelled?
- Done?
- createdBy
- createdAt
- updatedBy
- updatedAt

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

#### Lit-element pages

- home/login page
- add/edit appointment
- add/edit patient
- add/edit user
- view scheduling
- view activities log

## TODOS

- ~~create entities models, interfaces and errors~~
- ~~Add tables to db setup sql file - appointments, surgeries~~
- ~~write failing tests for each API (appointments, surgeries and patients)~~
- ~~write services for each entity~~
- ~~add postgres container for testing~~
- validation should be in services implementation
- ~~add role for doctors (doctors are users who can have schedules [appointments|surgeries])~~

