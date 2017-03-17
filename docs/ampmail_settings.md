# Setting email server in amp

Amp sends several kinds of emails, for instance to confirme an account creation or to reset a password. Amp doesn't have integrated email server and needs to use the services of an external email server.

This document explains how to set it.

## Default behavior

Without setting, amp uses the services of the email delivery site SendGrid: https://sendgrid.com/marketing/sendgrid-services

However amp needs an SendGrid api key to be able to use it.

To set the SendGrid apikey, create or update the file: `~/.config/amp/amplifier.yaml` and add this line inside:

EmailKey: xxxxx

where xxxxx is the SendGrid api key.

Then you have to rebuild the amp image and use it to deploy amp (at least use it to start amplifier).


## Set another external email server


To set any email server using SMTP protocol, add the following arguements in amplifier service:

- --email-server-address: the address of your email server
- --email-server-port: the SMTP port used by your email server
- --email-sender: the email which are going to be seen as the email sender

Amp handles automatically the fact that your server needs TLS or not.

To set your email sender password, create or update the file: `~/.config/amp/amplifier.yaml` and add this line inside:

EmailPwd: xxxxx

where xxxxx is your password

Then you have to rebuild the amp image and use it to deploy amp (at least use it to start amplifier).

