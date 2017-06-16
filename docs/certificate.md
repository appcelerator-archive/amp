# Certificate management for AMP

For local deployment, amp will generate a self signed certificate. You'll have to accept it when you connect the first time to the services on local.appcelerator.io.

For remote deployment, if you want to use a valid DNS domain, you can upload a certificate on the Swarm to enable TLS flows.

If your domain is `cloud.mydomain.net`, the certificate should be signed for `cloud.mydomain.net`, `dashboard.cloud.mydomain.net`, `gw.cloud.mydomain.net`, `examples.cloud.mydomain.net` and `service.cloud.domain.net`. If you can get a wildcard certificate, it's even better.

As an admin user of the platform, you can replace the certificate in the swarm. For that, prepare the secret as a pem file (includind the private key, the certificate, and the certificate chain), create a new Docker secret in the swarm, and update the amp_proxy service to mount it as `/run/secrets/certN.pem`, with N = 1 to 9.
```
amp secret create amp_certificate_$(date +%Y%m%d)
amp service update --secret-add amp_certificate_$(date +%Y%m%d) amp_proxy
amp service update --secret-rm amp_certificate_old_one amp_proxy
```
