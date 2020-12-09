# Dev

This folder contains files used to start up servers in development mode.

## Prerequisites

### A reverse proxy

It needs a pre configured basic reverse proxy to map a TLS enabled local
development DNS entry with the following paths:

- `/api` -> mapping to the backend server address/port
- `/` -> mapping to the frontend http server address/port

See [this blog post](https://carlosapgomes.me/post/localssl/) on how to
configure a local development nginx reverse proxy with HTTPS.

### Environment Variables

Create a shell script named `setEnv.sh` containing the code below (provide
a valid sendgrid api key and a valid email address):

```sh

#!/bin/sh
echo "Setting environment variables"
export PG_STR="postgres://sked:sked@localhost:54320/sked?sslmode=disable"
export SENDGRID_API_KEY="SG.xxxxxxxxxxxxxxxxxx"
export FROM_EMAIL="a_valid_email@address"

```

## How to run

From the current folder (`dev`) run the following command:

`source ./setEnv.sh && ./startDevMode.sh`

And wait until the backend server shows a message like this:

`Starting server on :9000 `

Now, go to the `../frontend/reactfrontend/` folder and start the
react development server with the command:

`yarn start`

Next, open a new browser tab and access the local development address
that was configured in the [Prerequisites](#prerequisites) section above.
