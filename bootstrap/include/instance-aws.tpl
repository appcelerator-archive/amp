{
  "Plugin": "instance-aws",
  "Properties": {
    "RunInstancesInput": {
      "ImageId": "{{ ref "/aws/amiid" }}",
      "InstanceType": "{{ ref "/aws/instancetype" }}",
      "KeyName": "{{ ref "/aws/keyname" }}",
      "SubnetId": "{{ ref "/aws/subnetid" }}",
      {{ if ref "/aws/instanceprofile" }}"IamInstanceProfile": {
        "Name": "{{ ref "/aws/instanceprofile" }}"
      },{{ end }}
      "SecurityGroupIds": [ "{{ ref "/aws/securitygroupid" }}" ]
    },
    "Tags": {
      "Name": "{{ ref "/aws/stackname" }}"
    }
  }
}
