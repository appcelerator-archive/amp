{
  "Plugin": "flavor-vanilla",
  "Properties": {
    "Init": [
      "docker run --rm  --network {{ ref "/amp/network" }} alpine sh -c 'nslookup $(hostname)'",
      "if [ $? -ne 0 ]; then echo 'Docker Swarm DNS check failed'; exit 1; fi",
      "exit 0"
    ]
  }
}
