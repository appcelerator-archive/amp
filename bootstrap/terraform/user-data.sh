#cloud-config
repo_update: true
repo_upgrade: security
packages:
  - ca-certificates
  - jq
  - git
  - curl
  - unzip
write_files:
  - path: /tmp/env.ikt
    content: |
      {{ global "/aws/stackname" "${instance_infrakit_group_suffix}" }}
runcmd:
  - wget -qO- https://get.docker.com/ | sh
  - usermod -G docker ubuntu
  - systemctl enable docker.service
  - systemctl start docker.service
  - curl ${infrakit_config_base_url}/bootstrap -o /usr/local/bin/bootstrap.sh
  - bash /usr/local/bin/bootstrap.sh -p terraform -t aws ${infrakit_config_base_url} -e /tmp/env.ikt
