# Go mailing service

Very simple service that sends emails via smtp. This is the template I use on my project that require that feature.

Requires the following environment variables:

- App port (that the server will listen on)
- Email from (email account that will send the emails)
- Email password (3rd party app email password)
- SMTP host
- SMTP port
- Recipients (optional)
- Subject (optional)

POST to `/api/send` with the request body containing the following variables:

- Recipients (optional)
- Subject (optional)
- Message
