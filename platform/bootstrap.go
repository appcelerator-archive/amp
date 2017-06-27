package main

import (
	"encoding/base64"
	"log"
	"strconv"

	"github.com/appcelerator/amp/pkg/docker"
	"github.com/docker/docker/api/types/swarm"

	"golang.org/x/net/context"
)

const (
	minimumApiVersion = 1.30
)

// initialize a secret. the data will be base64 encrypted
func initSecret(client *docker.Docker, name string, text []byte) (string, error) {
	data := base64.StdEncoding.EncodeToString([]byte(text))
	secretSpec := swarm.SecretSpec{Annotations: swarm.Annotations{Name: name}, Data: []byte(data)}
	resp, err := client.GetClient().SecretCreate(context.Background(), secretSpec)
	if err != nil {
		return "", err
	}
	return resp.ID, err
}

func main() {
	client := docker.NewClient(docker.DefaultURL, docker.DefaultVersion)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	c := client.GetClient()
	// check docker version
	version, err := c.ServerVersion(context.Background())
	log.Printf("Docker engine version %s\n", version.Version)
	apiVersion, err := strconv.ParseFloat(version.APIVersion, 32)
	if err != nil {
		log.Fatal(err)
	}
	if apiVersion < minimumApiVersion {
		log.Fatal("Docker engine doesn't meet the requirements (API Version)")
	}

	// secret initialization
	secretId, err := initSecret(client, "amplifier_yml", []byte("---"))
	if err != nil {
		log.Fatal("Failed to create the secret")
	}
	log.Printf("amplifier_yml secret: %s\n", secretId)
	secretId, err = initSecret(client, "alertmanager_yml", []byte("---"))
	if err != nil {
		log.Fatal("Failed to create the secret")
	}
	log.Printf("alertmanager_yml secret: %s\n", secretId)
	secretId, err = initSecret(client, "prometheus_alert-rules", []byte("---"))
	if err != nil {
		log.Fatal("Failed to create the secret")
	}
	log.Printf("prometheus_alert_rules secret: %s\n", secretId)
	secretId, err = initSecret(client, "certificate_atomiq", []byte("---"))
	if err != nil {
		log.Fatal("Failed to create the secret")
	}
	log.Printf("certificate_atomiq secret: %s\n", secretId)
	// does the swarm use labels
	// check size of cluster
	// check all labels are defined, if used
	// stack deployment
	// smoke tests
}
