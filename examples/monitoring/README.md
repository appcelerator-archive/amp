# Monitoring sample files

Extensive reference documentation is available on [prometheus.io].

## Validation of configuration files

### Prometheus alert rules

The configuration file can be tested with the following command:

    docker run --rm -v $PWD/examples/monitoring/prometheus_alerts.rules:/alerts.rules --entrypoint /bin/promtool prom/prometheus check rules /alerts.rules

## Deployment of new configuration

### Prometheus alert rules

To update the configuration, you first have to check the name of the Docker config used by the prometheus service, and list the existing Docker configs to be able to build a new unique name.

    amp -s <REMOTE_URL> service inspect amp_prometheus --format '{{ (index .Spec.TaskTemplate.ContainerSpec.Configs 0 ).ConfigName }}'
    amp -s <REMOTE_URL> config ls | grep prometheus
    amp -s <REMOTE_URL> config create prometheus_alerts_rules_CLUSTERNAME examples/monitoring/prometheus_alerts.rules
    amp -s <REMOTE_URL> service update --config-rm prometheus_alerts_rules --config-add source=prometheus_alerts_rules_CLUSTERNAME,target=/etc/prometheus/alerts.rules

### Alertmanager configuration

    amp -s <REMOTE_URL> secret ls | grep alertmanager
    amp -s <REMOTE_URL> service inspect amp_alertmanager --format '{{ (index .Spec.TaskTemplate.ContainerSpec.Secrets 0 ).SecretName }}'
    amp -s <REMOTE_URL> secret create alertmanager_yml_CLUSTERNAME examples/monitoring/prometheus_alerts.rules
    amp -s <REMOTE_URL> service update --config-rm prometheus_alerts_rules --config-add source=prometheus_alerts_rules_CLUSTERNAME,target=/etc/prometheus/alerts.rules
