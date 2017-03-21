#cloud-config
repo_update: true
repo_upgrade: security
packages:
  - ca-certificates
  - jq
  - git
  - curl
  - unzip
  - make
write_files:
  - path: /root/.config/infrakit/infrakit/env.ikt
    content: |
      {{/* Global variables */}}
      {{ global "/terraform/region" "${region}" }}
      {{ global "/terraform/stackname" "${name}" }}
      {{ global "/terraform/vpcid" "${vpc_id}" }}
      {{ global "/terraform/subnetid" "${subnet_id}" }}
      {{ global "/terraform/securitygroupid" "${security_group_id}" }}
      {{ global "/terraform/amiid" "${ami}" }}
      {{ global "/terraform/instancetype" "${instance_type}" }}
      {{ global "/terraform/instanceprofile" "${cluster_instance_profile}" }}
      {{ global "/terraform/keyname" "${key_name}" }}
      {{ global "/script/baseurl" "${infrakit_config_base_url}" }}
      {{ global "/docker/aufs/size" "${aufs_volume_size}" }}
runcmd:
  - wget -qO- https://get.docker.com/ | sh
  - usermod -G docker ubuntu
  - systemctl enable docker.service
  - systemctl start docker.service
  - curl ${infrakit_config_base_url}/bootstrap -o /usr/local/bin/bootstrap.sh
  - bash /usr/local/bin/bootstrap.sh -p terraform ${infrakit_config_base_url}
