# Dev

This folder contains config files used in development mode:

- postgres docker files
- start up scripts

To run the server, make sure to provide the following command line
parameters:

- addr (address and port to listen to)
- pgstr (connection string to postgres)
- sgKey (Sendgrid API key)
- fromAddress

And a local HTTPS reverse proxy
