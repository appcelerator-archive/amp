# Setting sms sender in amp

Amp sends several kinds of sms, for instance to double check an account creation or to reset a password. Amp use the services of an external sms server.

This document explains how to set it.

## Default behavior

amp uses the services of the sms delivery site twilio: https://api.twilio.com

amp needs an twilio api key and an account id to be able to use it.

To set then, create or update the file: `~/.config/amp/amplifier.yaml` and add this line inside:

SmsKey: xxxxx
SmsAccountId: yyyyy

where xxxxx is the twilio api key and yyyyy the account id

It's also possible to set the sms sender:

SmsSender: xxxxx
where xxxxx is the number which appear as the sender for all amp sms


Then you have to rebuild the amp image and use it to deploy amp (at least use it to start amplifier).
