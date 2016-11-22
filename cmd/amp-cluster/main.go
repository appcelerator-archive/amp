package main

// build vars
var (
	Version string
	Build   string
	client  = &clusterClient{}
)

func main() {
	client.cli()
}
