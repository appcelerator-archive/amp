##### Secrets

`amp cluster create` uses a docker secret named `amplifier_yml` for amplifier configuration.

If the secret is not present before the invocation of `amp cluster create`, it will be automatically generated with sensible values for the following keys:
- `JWTSecretKey`: A secret key of 128 random characters will be generated.
- `SUPassword`: A super user password of 32 characters will be generated and displayed during the execution of the command.

If the secret is already created, it will be used as is without any modifications.
