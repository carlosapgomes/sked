# Sked(duler) [![time tracker](https://wakatime.com/badge/github/carlosapgomes/sked.svg)](https://wakatime.com/badge/github/carlosapgomes/sked)

This software intends to substitute a paper-based planner that I used to use
for keeping track of patients, appointments, and scheduled surgeries in my
medical practice.

## Background

As a context, I work at a day hospital that provides medical services to the
public health care system in my country. We need to triage those patients
before surgery and to follow up for some time depending on the procedure,
before returning them to the referring clinic.

We receive a fixed number of patients per week for a triage appointment, but
we also have to schedule all follow-up consultations in a way that does not
overload our capacity. Besides that, I need to keep track of my scheduled
surgeries to balance the surgery types for each day to optimize for time
to patient discharge.

All that information is kept in two paper-based planners. Scheduled surgeries
and some appointments are registered in a planner that I carry along with me
and the rest in another one that is with the day hospital secretary. The
main problem that I want to solve is to synchronize those pieces of
information and let them available to both of us.

## Installation

There is an ansible playbook in the build folder that installs and configures
the system in a Debian server (I did not test it with other distros).

It requires access to the server with a public key authentication.

Create an inventory file with all the required values.

If all you want is just test the app, just go to the `dev` folder and follow
the instructions in the
[README](https://github.com/carlosapgomes/sked/blob/master/dev/README.md) file.

## Usage

It is designed to be a web-based application.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
