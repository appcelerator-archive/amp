{
  "Plugin": "flavor-vanilla",
  "Properties": {
    "Init": [
      "# create an overlay network",
      "docker network inspect {{ ref "/amp/network" }} 2>&1 | grep -q 'No such network' && \\",
      "  docker network create -d overlay --attachable {{ ref "/amp/network" }}",
      "exit 0"
    ]
  }
}
