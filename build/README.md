# Build & Deploy

Scripts and configurations for building and deploying the system.

## Ansible

Configure a VM to be used by Ansible
Create an inventory file in the current `build` folder, for example
`ansible_host`, and add a new entry in the `.gitignore` file on this
project root folder so that it is not pushed for the remote git repository.

In the inventory file, provide the following environment variables:

- PG*PASSWD - the password to be used by Postgres
  * PG_LC - locale to configure Postgres database
- SG_API_KEY - Sendgrid key
- FROM_EMAIL - email address to be used as a sender by the mailer
- FROM_EMAIL - defaults to "Sked Manager"
- TZ - backend timezone
