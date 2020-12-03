# Build & Deploy

Scripts and configurations for building and deploying the system.

## Ansible

Configure a VM to be used by Ansible
Create an inventory file in the current `build` folder, for example
`ansible_inv.yaml`, and add a new entry in the `.gitignore` file on this
project root folder so that it is not pushed for the remote git repository
(see `ansible_inv_sample.yaml`).

In the inventory file, provide the following environment variables:

- sys_timezone - backend timezone
- PG_PASSWD - the password to be used by Postgres
- PG_LC - locale to configure Postgres database
- SG_API_KEY - Sendgrid key
- FROM_EMAIL - email address to be used as a sender by the mailer
- FROM_NAME - defaults to "Sked Manager"

From the current folder (`build`) run:

`ansible-playbook -i ansible_inv.yml provisioning/playbook.yml`
