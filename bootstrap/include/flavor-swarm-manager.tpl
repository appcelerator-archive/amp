{
  "Plugin": "flavor-swarm/manager",
  "Properties": {
    "InitScriptTemplateURL": "{{ ref "/script/baseurl" }}/manager-init.tpl",
    "SwarmJoinIP": "{{ ref "managerSwarmJoinIP" }}",
    "Docker" : {
      {{ if ref "/certificate/ca/service" }}"Host" : "{{ ref "managerDockerHostTLS" }}",
      "TLS" : {
        "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
        "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
        "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
        "InsecureSkipVerify": false
      }
      {{ else }}"Host" : "{{ ref "managerDockerHost" }}" {{ end }}
    }
  }
}
