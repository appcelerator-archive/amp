{
  "Plugin": "flavor-swarm/worker",
  "Properties": {
    "InitScriptTemplateURL": "{{ ref "/script/baseurl" }}/worker-init.tpl",
    "SwarmJoinIP": "{{ref "workerSwarmJoinIP" }}",
    "Docker" : {
      {{ if ref "/certificate/ca/service" }}"Host" : "{{ ref "workerDockerHostTLS" }}",
      "TLS" : {
        "CAFile": "{{ ref "/docker/remoteapi/cafile" }}",
        "CertFile": "{{ ref "/docker/remoteapi/certfile" }}",
        "KeyFile": "{{ ref "/docker/remoteapi/keyfile" }}",
        "InsecureSkipVerify": false
      }
      {{ else }}"Host" : "{{ ref "workerDockerHost" }}" {{ end }}
    }
  }
}
