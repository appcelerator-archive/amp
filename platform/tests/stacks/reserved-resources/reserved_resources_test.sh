#!/usr/bin/env bash

test_reserved_secret_amplifier() {
  amp -k stack up -c platform/tests/stacks/reserved-resources/secret.amplifier.yml
  return $(($? ^ 1))
}

test_reserved_secret_certificate() {
  amp -k stack up -c platform/tests/stacks/reserved-resources/secret.certificate.yml
  return $(($? ^ 1))
}

test_reserved_label_io_amp_role() {
  amp -k stack up -c platform/tests/stacks/reserved-resources/label.io.amp.role.yml
  return $(($? ^ 1))
}
