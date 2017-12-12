# AWS cluster creation

This page outlines the steps to set up an AMP cluster on AWS.

## Why do you need an AMP cluster on AWS?

When you application is ready or when you want your team to share a common cluster, it is easy to spin up an AMP cluster on AWS. 

## Prerequisites

The prerequisites for creating an AWS cluster on AMP is to have:
 
### AWS Account
 
You will need an AWS account with adequate IAM rights to deploy Cloud Formation stacks. 

If you don't already have one, sign up [here](https://portal.aws.amazon.com/billing/signup#/start). 

To learn more about managing IAM roles, see the official documentation from AWS [here](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html?icmpid=docs_iam_console).
 
### AWS Key Pair 
 
You will need an AWS key pair specific to the region you choose. 

For creating key pairs on AWS, see the official documentation from AWS [here](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html).

### AWS Access Key

You will need to provide your `AWS Access Key ID` and `AWS Secret Access Key` in order to deploy an AMP cluster.

You can find out more about AWS access keys [here](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey).

> NOTE: If you already have the AWS CLI installed on your system, make sure it is configured since the AMP CLI will use this information for cluster creation. See more details [here](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html). 

## Checklist

Before deploying your cluster, you will need the following elements:

* Pick a region, for instance `us-west-2`
* AWS Key pair name in the region you chose
* AWS Access Key ID
* AWS Secret Access Key
* Stack name

> NOTE: If you already have the AWS CLI installed, the AMP CLI will use the `Region`, `AWS Access Key ID` and `AWS Secret Access Key` from it.

## Cluster Deployment

If you don't have the AWS CLI installed, enter the following:
```
$ export REGION=us-west-2
$ export KEY_NAME=user-keypair
$ export ACCESS_KEY_ID=xxxxx
$ export SECRET_ACCESS_KEY=xxxxx
$ export STACK_NAME=amp-test

$ amp cluster create --provider aws --aws-region $REGION --aws-parameter KeyName=$KEY_NAME --aws-access-key-id $ACCESS_KEY_ID --aws-secret-access-key $SECRET_ACCESS_KEY --aws-stackname $STACK_NAME --aws-sync
```

If you have the AWS CLI installed, enter the following:
```
$ export KEY_NAME=user-keypair
$ export STACK_NAME=amp-test

$ amp cluster create --provider aws --aws-parameter KeyName=$KEY_NAME --aws-stackname $STACK_NAME --aws-sync
```

The AMP cluster deployment on AWS takes roughly 10 minutes to complete. On success, the output looks something like this:
```
Fri Dec  8 01:16:53 UTC 2017    amp-test                     CREATE_IN_PROGRESS (User Initiated)
Fri Dec  8 01:17:16 UTC 2017    Vpc                          CREATE_COMPLETE
Fri Dec  8 01:19:34 UTC 2017    ManagerAutoScalingGroup      CREATE_COMPLETE
Fri Dec  8 01:20:56 UTC 2017    ManagerWaitCondition         CREATE_COMPLETE
Fri Dec  8 01:21:02 UTC 2017    CoreWorkerAutoScalingGroup   CREATE_COMPLETE
Fri Dec  8 01:21:02 UTC 2017    UserWorkerAutoScalingGroup   CREATE_COMPLETE
Fri Dec  8 01:21:58 UTC 2017    CoreWaitCondition            CREATE_COMPLETE
Fri Dec  8 01:21:59 UTC 2017    UserWaitCondition            CREATE_COMPLETE
Fri Dec  8 01:23:50 UTC 2017    ApplicationWaitCondition     CREATE_COMPLETE
Fri Dec  8 01:23:54 UTC 2017    amp-test                     CREATE_COMPLETE
-------------------------------------------------------------------------------
VPC ID                                     | vpc-a92ec8d0
NFSv4 Endpoint                             | disabled
URL for cluster health dashboard           | amp-test-ManagerEx-8IKHVRJOUT98-636105687.us-west-2.elb.amazonaws.com:8080
internal endpoint for the registry service | disabled
public facing endpoint for the cluster     | amp-test-ManagerEx-8IKHVRJOUT98-636105687.us-west-2.elb.amazonaws.com
```

## What's next?

You can now deploy a stackfile on your newly created local cluster. Please follow the instructions listed [here](stackdeploy.md).

## Customization

The cluster has been created with the default configuration. You'll need to update it according to your convenience.

## Domain Certificate

You should generate a certificate valid for the domain that will be served by the cluster. Ideally it can be a wildcard certificate, but it should at least include the virtual hosts amplifier, dashboard, kibana and alerts.

The default certificate is stored in the secret `certificate_amp`. It should include the private key, the certificate and the full certificate chain (in this order). Create a new secret with your new certificate:

```
$ cat certificate.pem | docker secret create certificate_amp_custom -
```

Update the services that use it:

```
$ docker service update --secret-rm certificate_amp --secret-add source=certificate_amp_custom,target=/run/secrets/cert0.pem amp_amplifier
$ docker service update --secret-rm certificate_amp --secret-add source=certificate_amp_custom,target=/run/secrets/cert0.pem amp_proxy
```

Check that the service is stabilizing with the new configuration

```
$ docker service ps amp_amplifier
$ docker service ps amp_proxy
```

## Secrets

`amp cluster create` uses a docker secret named `amplifier_yml` for amplifier configuration.

If the secret is not present before the invocation of `amp cluster create`, it will be automatically generated with sensible values for the following keys:
- `JWTSecretKey`: A secret key of 128 random characters will be generated.
- `SUPassword`: A super user password of 32 characters will be generated and displayed during the execution of the command.

If the secret is already created, it will be used as is without any modifications.

## Amplifier configuration

The secret `amplifier_yml` contains the configuration for the user registration.
To enable it, you should get a Sendgrid key (you can have a 30 days trial key [here](https://app.sendgrid.com/signup?id=71713987-9f01-4dea-b3d4-8d0bcd9d53ed).

```
$ cat > amplifier.yml << EOF
EmailKey: SENDGRID_KEY
SUPassword: SUPER_USER_PASSWORD
JWTSecretKey: JWT_RANDOM_STRING
EOF
$ cat amplifier.yml | docker secret create amplifier_yml_custom -
```

Update the service:

```
$ docker service update --secret-rm amplifier_yml --secret-add source=amplifier_yml_custom,target=/run/secrets/amplifier.yml amp_amplifier
```

Check that the service is stabilizing with the new configuration.

## Alerting configuration

You can configure where notifications will be sent. The default (no notification configured) is in the secret `alertmanager_yml`.

Create a new file with your configuration:
```
$ cat > alertmanager.yml << EOF
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
$ cat alertmanager.yml | docker secret create alertmanager_yml_custom -
```

Update the service:

```
$ docker service update --secret-rm alertmanager_yml --secret-add source=alertmanager_yml_custom,target=/run/secrets/alertmanager.yml amp_alertmanager
```

Check that the service is stabilizing with the new configuration.

Now you can also configure the alerts definitions. By default, the Docker config `prometheus_alerts_rules` contains empty rules.

Create a new file with the configuration. You can find an example in the amp repo in `examples/monitoring/prometheus_alerts.rules`.

```
$ cat prometheus_alerts.rules | docker config create prometheus_alerts_rules_custom -
```

Update the service

```
$ docker service update --config-rm prometheus_alerts_rules --config-add source=prometheus_alerts_rules_custom,target=/etc/prometheus/alerts.rules amp_prometheus
```

Check that the service is stabilizing with the new configuration.

## Update the domain name

You should create or update the domain name that you use to create the certificate. Use a service such as route53, dnsimple or similar to create a CNAME pointing to the DNS name from the output of the cluster creation. If you don't have access to it anymore, you can find it in the output of the cloudformation stack on the AWS console.

Wait for the DNS to be updated.

## Login and validate the configuration

Use the CLI from your workstation (not from the cluster anymore) and check that the connection is done and that the certificate is valid:

```
$ amp -s DOMAIN_NAME user ls
```

Then signup a new user

```
$ amp -s DOMAIN_NAME user signup --name LOGIN --email EMAIL
```

You should receive an email with a command line for the validation of the account.

Note that if it is a test cluster, you can skip certificate creation step and use -k (--insecure) AMP CLI argument.

Also you can grab public endpoint (ELB) DNS name from Cloudformation -> Stacks - Outputs -> DNSTarget key.

```
$ DOMAIN_NAME=nderzhak-ManagerE-1FQNBP6KF3O24-818587690.us-west-2.elb.amazonaws.com; ./amp -s $DOMAIN_NAME -k user ls
[nderzhak-ManagerE-1FQNBP6KF3O24-818587690.us-west-2.elb.amazonaws.com:50101]
USERNAME   EMAIL   CREATED ON
su                 29 Nov 17 14:41
```

## Cluster Cleanup

You can remove the old secrets and config now that the cluster is configured and running.

```
$ docker secret rm amplifier_yml
$ docker secret rm alertmanager_yml
$ docker secret rm certificate_amp
$ docker config rm prometheus_alerts_rules
```

## Cluster Teardown 

If you no longer use the deployed cluster, it can be removed by running the following command:

```
$ amp cluster rm --provider aws --aws-stackname $STACK_NAME --aws-sync
```

The cluster teardown on AWS takes about 10 minutes to complete. On success, the output looks something like this:
```
Fri Dec  8 01:17:16 UTC 2017    Vpc                          CREATE_COMPLETE
Fri Dec  8 01:19:34 UTC 2017    ManagerAutoScalingGroup      CREATE_COMPLETE
Fri Dec  8 01:20:56 UTC 2017    ManagerWaitCondition         CREATE_COMPLETE
Fri Dec  8 01:21:02 UTC 2017    UserWorkerAutoScalingGroup   CREATE_COMPLETE
Fri Dec  8 01:21:02 UTC 2017    CoreWorkerAutoScalingGroup   CREATE_COMPLETE
Fri Dec  8 01:21:58 UTC 2017    CoreWaitCondition            CREATE_COMPLETE
Fri Dec  8 01:21:59 UTC 2017    UserWaitCondition            CREATE_COMPLETE
Fri Dec  8 01:23:50 UTC 2017    ApplicationWaitCondition     CREATE_COMPLETE
Fri Dec  8 01:23:54 UTC 2017    amp-test                     CREATE_COMPLETE
Fri Dec  8 01:26:12 UTC 2017    amp-test                     DELETE_IN_PROGRESS (User Initiated)
Fri Dec  8 01:26:14 UTC 2017    ApplicationWaitCondition     DELETE_IN_PROGRESS
Fri Dec  8 01:26:15 UTC 2017    UserWaitCondition            DELETE_IN_PROGRESS
Fri Dec  8 01:26:15 UTC 2017    CoreWaitCondition            DELETE_IN_PROGRESS
Fri Dec  8 01:26:17 UTC 2017    UserWorkerAutoScalingGroup   DELETE_IN_PROGRESS
Fri Dec  8 01:26:17 UTC 2017    CoreWorkerAutoScalingGroup   DELETE_IN_PROGRESS
Fri Dec  8 01:29:54 UTC 2017    ManagerWaitCondition         DELETE_IN_PROGRESS
Fri Dec  8 01:29:56 UTC 2017    ManagerAutoScalingGroup      DELETE_IN_PROGRESS
Fri Dec  8 01:35:11 UTC 2017    Vpc                          DELETE_IN_PROGRESS
Fri Dec  8 01:35:30 UTC 2017    amp-test                     DELETE_COMPLETE
```
