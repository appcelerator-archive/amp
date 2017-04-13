# Setting email server in amp

Amp sends several kinds of emails, for instance to confirm an account creation or to reset a password. Amp doesn't have integrated email server and needs to use the services of an external email server.

This document explains how to set it.

## Default behavior

amp uses the services of the email delivery site SendGrid: https://sendgrid.com/marketing/sendgrid-services

amp needs an SendGrid api key to be able to use it.

To set the SendGrid apikey, create or update the file: `~/.config/amp/amplifier.yml` and add this line inside:

EmailKey: xxxxx

where xxxxx is the SendGrid api key.

It's also possible to set the email sender:

EmailSender: xxxxx
where xxxxx is the email which appear as the sender for all amp emails


Then you have to rebuild the amp image and use it to deploy amp (at least use it to start amplifier).
