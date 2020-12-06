# Build & Deploy

Scripts and configurations for building and deploying the system.

## Ansible

Configure a VM to be used by Ansible (I suggest using a Debian VM, because
some vars are specific to this distribution - `ansible_python_interpreter `
and `postgresql_python_library` are used by `geerlingguy.postgresql` role).

Create an inventory file in the current `build` folder, for example
`ansible_inv.yml`, and add a new entry in the `.gitignore` file on this
project root folder so that it is not pushed for the remote git repository
(see `ansible_inv_sample.yml`).

In the inventory file, provide the following variables:

- ansible_host: host.ip.address
- ansible_port: host_ssh_port
- ansible_user: host_user_name
- ansible_ssh_private_key_file: path/to/ssh/key/file
- sys_timezone - backend timezone
- admin_email - host admin email address
- pg_passwd - the password for user 'sked' in Postgres
- pg_loc - locale to configure Postgres database
- sg_api_key - Sendgrid key
- from_email - email address to be used as a sender by the mailer
- from_name - defaults to "Sked Manager"

From the current folder (`build`) run:

`ansible-playbook -i ansible_inv.yml provisioning/playbook.yml`
