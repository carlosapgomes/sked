# Dev

This folder contains config files used in development mode:

- postgres docker files
- start up scripts

To run the server, make sure to provide at least the following command line
parameters:

- pgstr (connection string to postgres)
- sgKey (Sendgrid API key)
- fromAddress

And a local HTTPS reverse proxy
