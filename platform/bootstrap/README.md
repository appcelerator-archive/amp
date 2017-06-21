# Deploy a Docker Swarm cluster with InfraKit

## Bootstrap

First, the cluster should be bootstrapped to get InfraKit running and ready to deploy the Swarm cluster.

It can be either:

#### Local swarm manager node

You only need Docker 17.03+ installed on your machine, and run the CLI to deploy AMP on it. It will initialize the swarm if it's not already done.

```amp -s localhost cluster create --registration=none --notifications=false```

#### AWS

##### Externally managed

Once you have a Swarm cluster with a domain name, ssh to one manager, download the amp CLI, and execute:

```amp -s localhost cluster create --domain cloud.domain.na.me --secrets-dir /absolute/path/to/secrets```

the secrets dir should contain:
- a domain.na.me.pem file with the certificate for the domain. It can be a wildcard certificate (*.cloud.domain.na.me) or contain a few vhosts (cloud.domain.na.me, gw.cloud.domain.na.me, alerts.cloud.domain.na.me, dashboard.cloud.domain.na.me, kibana.cloud.domain.na.me) 
- a amplifier.yml file with a SendGrid `EmailKey`, a `SuPassword` and a `JWTSecretKey`

##### Managed by InfraKit

The CLI will create a VPC, subnet, internet gateway and the minimum required to build EC2 instances.
Make sure you have AWS credentials ready ($HOME/.aws/credentials).

```amp -s localhost cluster create --provider aws --domain cloud.domain.na.me --secrets-dir /absolute/path/to/secrets```

the secrets dir should contain:
- a domain.na.me.pem file with the certificate for the domain. It can be a wildcard certificate (*.cloud.domain.na.me) or contain a few vhosts (cloud.domain.na.me, gw.cloud.domain.na.me, alerts.cloud.domain.na.me, dashboard.cloud.domain.na.me, kibana.cloud.domain.na.me) 
- a amplifier.yml file with a SendGrid `EmailKey`, a `SuPassword` and a `JWTSecretKey`
- a aws-parameter.json for the customization of your deployment, it should contain a content similar to:
```
[
  {
    "ParameterKey": "KeyName",
    "ParameterValue": "KEYPAIR_NAME"
  },
  {
    "ParameterKey": "ManagerSize",
    "ParameterValue": "3"
  },
  {
    "ParameterKey": "ManagerInstanceType",
    "ParameterValue": "t2.medium"
  },
  {
    "ParameterKey": "InfraKitConfigurationBaseURL",
    "ParameterValue": "https://raw.githubusercontent.com/appcelerator/amp/master/platform/bootstrap"
  }
]
```

#### DigitalOcean, GCP, Azure

See above the `Externally managed` section for AWS.
