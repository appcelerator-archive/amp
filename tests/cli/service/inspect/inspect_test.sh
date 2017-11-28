#!/bin/bash

test_inspect() {
  amp -k service inspect pinger_pinger 2>/dev/null | grep -q "pinger"
}

test_inspect_format() {
  secret_name=amplifier_yml
  secrets=$(amp -k service inspect amp_amplifier --format '{{ range .Spec.TaskTemplate.ContainerSpec.Secrets}} {{ .SecretName }} {{ end }}') || return 1
 echo $secrets | grep -q $secret_name
}
