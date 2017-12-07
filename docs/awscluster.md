#### Creating a cluster on AWS

To target AWS, you should use the --provider aws option, refer to the help for details on aws options: `amp cluster create -h`.

##### Prerequisites

the AMP CLI can be used to create a new AMP cluster on AWS. The prerequisite is to have access to an AWS account with enough IAM rights to deploy a stack (with IAM policies creation), and to have the credentials defined in the default profile of `$HOME/.aws/credentials`. If you don't, download the [AWS cli](http://docs.aws.amazon.com/cli/latest/userguide/installing.html), and run `aws configure`. You can check with `aws configure list`.

You can download CLI here: https://github.com/appcelerator/amp/releases

##### Cluster creation

Make sure to define the variables with the proper region (i.e. us-west-2), stack name (i.e. amp-YOURNAME) and keypair name. This keypair should already exist on AWS (you can list them in the AWS EC2 page for the chosen region).

```
amp cluster create --provider aws --aws-region $REGION --aws-stackname $STACK_NAME --aws-parameter KeyName=$KEY_NAME
```

The CLI will give back control once the cluster is fully deployed, which should take 10 min. It will display some information, including the ELB dns name for the manager nodes (which can be used to configure a CNAME with the domain served by this cluster).

##### Customization

The cluster has been created with the default configuration. You'll need to update it for your convenience.

For 0.14, the customization should be done on one manager node of the cluster, so you should be able to ssh to it with the key you specified at cluster creation, with user `ubuntu`.

##### Certificate

You should generate a certificate valid for the domain that will be served by the cluster. Ideally it can be a wildcard certificate, but it should at least include the virtual hosts amplifier, dashboard, kibana and alerts.

The default certificate is stored in the secret `certificate_amp`. It should include the private key, the certificate and the full certificate chain (in this order). Create a new secret with your new certificate:

```
cat certificate.pem | docker secret create certificate_amp_custom -
```

Update the services that use it:

```
docker service update --secret-rm certificate_amp --secret-add source=certificate_amp_custom,target=/run/secrets/cert0.pem amp_amplifier
docker service update --secret-rm certificate_amp --secret-add source=certificate_amp_custom,target=/run/secrets/cert0.pem amp_proxy
```

Check that the service is stabilizing with the new configuration

```
docker service ps amp_amplifier
docker service ps amp_proxy
```

##### Amplifier configuration

The secret `amplifier_yml` contains the configuration for the user registration.
To enable it, you should get a Sendgrid key (you can have a 30 days trial key [here](https://app.sendgrid.com/signup?id=71713987-9f01-4dea-b3d4-8d0bcd9d53ed)).

```
cat > amplifier.yml << EOF
EmailKey: SENDGRID_KEY
SUPassword: SUPER_USER_PASSWORD
JWTSecretKey: JWT_RANDOM_STRING
EOF
cat amplifier.yml | docker secret create amplifier_yml_custom -
```

Update the service:

```
docker service update --secret-rm amplifier_yml --secret-add source=amplifier_yml_custom,target=/run/secrets/amplifier.yml amp_amplifier
```

Check that the service is stabilizing with the new configuration.

##### Alerting configuration

You can configure where notifications will be sent. The default (no notification configured) is in the secret `alertmanager_yml`.

Create a new file with your configuration:
```
cat > alertmanager.yml << EOF
---
global:
  slack_api_url: "https://hooks.slack.com/services/YOUR_WEB_HOOK"
templates:
- '/etc/alertmanager/template/*.tmpl'
route:
  receiver: 'slack-receiver'
  repeat_interval: 5m
  routes:
  - receiver: 'slack-receiver'
receivers:
- name: "slack-receiver"
  slack_configs:
  - channel: "@CHANNEL_NAME"
    username: "Alertmanager"
    title: '*DOMAIN_NAME* {{ range .Alerts }}{{ .Annotations.summary }} {{ end }}'
    text: '{{ range .Alerts }}{{ .Annotations.description }} {{ end }}'
    send_resolved: true
EOF
cat alertmanager.yml | docker secret create alertmanager_yml_custom -
```

Update the service:

```
docker service update --secret-rm alertmanager_yml --secret-add source=alertmanager_yml_custom,target=/run/secrets/alertmanager.yml amp_alertmanager
```

Check that the service is stabilizing with the new configuration.

Now you can also configure the alerts definitions. By default, the Docker config `prometheus_alerts_rules` contains empty rules.

Create a new file with the configuration. You can find an example in the amp repo in `examples/monitoring/prometheus_alerts.rules`.

```
cat prometheus_alerts.rules | docker config create prometheus_alerts_rules_custom -
```

Update the service

```
docker service update --config-rm prometheus_alerts_rules --config-add source=prometheus_alerts_rules_custom,target=/etc/prometheus/alerts.rules amp_prometheus
```

Check that the service is stabilizing with the new configuration.

##### Update the domain name

You should create or update the domain name that you use to create the certificate. Use a service such as route53, dnsimple or similar to create a CNAME pointing to the DNS name from the output of the cluster creation. If you don't have access to it anymore, you can find it in the output of the cloudformation stack on the AWS console.

Wait for the DNS to be updated.

##### Login and validate the configuration

Use the CLI from your workstation (not from the cluster anymore) and check that the connection is done and that the certificate is valid:

```
amp -s DOMAIN_NAME user ls
```

Then signup a new user

```
amp -s DOMAIN_NAME user signup --name LOGIN --email EMAIL
```

You should receive an email with a command line for the validation of the account.

Note. If it is test cluster you can skip certificate creation step and use -k (--insecure) amp CLI argument.

Also you can grab public endpoint (ELB) DNS name from Cloudformation -> Stacks - Outputs -> DNSTarget key.

```
nderzhakmba-2:doc-test nikolaiderzhak$ DOMAIN_NAME=nderzhak-ManagerE-1FQNBP6KF3O24-818587690.us-west-2.elb.amazonaws.com; ./amp -s $DOMAIN_NAME -k user ls
[nderzhak-ManagerE-1FQNBP6KF3O24-818587690.us-west-2.elb.amazonaws.com:50101]
USERNAME   EMAIL   CREATED ON
su                 29 Nov 17 14:41
```

##### Cleanup

You can remove the old secrets and config now that the cluster is configured and running.

```
docker secret rm amplifier_yml
docker secret rm alertmanager_yml
docker secret rm certificate_amp
docker config rm prometheus_alerts_rules
```

##### Teardown 

Below is example of command to sunset some cluster.

```
REGION=us-west-2; STACK_NAME=nderzhak-test016; ./amp cluster rm --provider aws --aws-stackname $STACK_NAME --aws-region $REGION
```
