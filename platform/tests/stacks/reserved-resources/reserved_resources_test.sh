#!/usr/bin/env bash

test_setup() {
  amp="amp -s localhost"
  $amp user signup --name user1 --password password --email email@user1.amp
}

test_reserved_secret_amplifier() {
  $amp stack up -c platform/tests/stacks/reserved-resources/secret.amplifier.yml
  return $(($? ^ 1))
}

test_reserved_secret_certificate() {
  $amp stack up -c platform/tests/stacks/reserved-resources/secret.certificate.yml
  return $(($? ^ 1))
}

test_reserved_label_io_amp_role() {
  $amp stack up -c platform/tests/stacks/reserved-resources/label.io.amp.role.yml
  return $(($? ^ 1))
}
