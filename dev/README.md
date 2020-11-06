# Dev

This folder contains config files used in development mode:

- postgres docker image startup script
- db configuration scripts

To run the compiled server binary, make sure to provide at least the
following command line parameters:

- pgstr (connection string to postgres)
- sgKey (Sendgrid API key)
- from email address

And set up a local dns entry and an HTTPS reverse proxy to pass connections
to `localhost:9000` as described
[here](https://carlosapgomes.me/post/localssl/).

## Example

```sh

# start the database
./startPg.sh

# start sked backend
skedBackend -pgstr "postgres://user:password@localhost/sked?sslmode=disable" \\
            -sgKey SG.xxxxxxx \\
            -from "manager@domain.sked"
```

For this dev setup, use
`postgres://sked:sked@localhost:54320/sked?sslmode=disable` as postgres
connection string.

Another option is to provide the following environment variables and just
run the binary:

- PG_STR
- SENDGRID_API_KEY
- FROM_EMAIL

The binary looks for a folder in the current directory called `templates` which
contain all the templates for the HTML pages and emails.
