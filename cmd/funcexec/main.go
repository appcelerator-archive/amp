package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/appcelerator/amp/data/functions"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"golang.org/x/net/context"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// natsStreaming is the nats streaming client
	natsStreaming *ns.NatsStreaming

	// docker is the Docker client
	dock *docker.Docker
)

const (
	name                = "funchttp"
	connectionTimeout   = 5 * time.Second
	functionExecTimeout = time.Minute
)

func main() {
	log.Printf("%s (version: %s, build: %s)\n", name, Version, Build)

	// Nats
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname:", err)
	}
	natsStreaming = ns.NewClient(ns.DefaultURL, ns.ClusterID, name+"-"+hostname, connectionTimeout)
	if err := natsStreaming.Connect(); err != nil {
		log.Fatalln(err)
	}

	// NATS, subscribe to function topic
	log.Println("Subscribing to subject:", ns.FunctionSubject)
	_, err = natsStreaming.GetClient().QueueSubscribe(ns.FunctionSubject, ns.FunctionQueue, messageHandler, stan.DeliverAllAvailable())
	if err != nil {
		natsStreaming.Close()
		log.Fatalln("Unable to subscribe to subject", err)
	}
	log.Println("Subscribed to subject:", ns.FunctionSubject)

	// Docker
	dock = docker.NewClient(docker.DefaultURL, docker.DefaultVersion)
	log.Printf("Connecting to Docker API at %s version API: %s\n", docker.DefaultURL, docker.DefaultVersion)
	if err := dock.Connect(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to Docker API at", docker.DefaultURL)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Println("\nReceived an interrupt, unsubscribing and closing connection...")
			natsStreaming.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func messageHandler(msg *stan.Msg) {
	go processMessage(msg)
}

func processMessage(msg *stan.Msg) {
	// Parse function call message
	functionCall := functions.FunctionCall{}
	err := proto.Unmarshal(msg.Data, &functionCall)
	if err != nil {
		log.Println("Error unmarshalling function call:", err)
		return
	}
	log.Println("Function call received, call id:", functionCall.CallID)

	ctx := context.Background()

	// Create container
	config := &container.Config{
		Image:        functionCall.Function.Image,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		StdinOnce:    true,
	}
	hostConfig := &container.HostConfig{}
	networkingConfig := &network.NetworkingConfig{}
	container, err := dock.ContainerCreate(ctx, config, hostConfig, networkingConfig, stringid.GenerateNonCryptoID())
	if err != nil {
		log.Println("Error creating container:", err)
		return
	}
	log.Println("Container created:", container.ID)

	// Start
	if err = containerStart(ctx, container.ID); err != nil {
		log.Println("error starting container:", err)
		return
	}
	log.Println("Container started")

	// Attach container streams
	attachment, err := containerAttach(ctx, container.ID)
	if err != nil {
		log.Println("error attaching container:", err)
		return
	}
	defer attachment.Close()
	log.Println("Container attached")

	// Standard input
	stdIn := bufio.NewReader(bytes.NewReader(functionCall.Input))

	// Standard output
	var stdOutBuffer bytes.Buffer
	stdOut := bufio.NewWriter(&stdOutBuffer)

	// Standard error
	var stdErrorBuffer bytes.Buffer
	stdErr := bufio.NewWriter(&stdErrorBuffer)

	// Handle standard streams
	streamCtx, cancel := context.WithTimeout(ctx, functionExecTimeout)
	defer cancel()
	handleStreams(streamCtx, stdIn, stdOut, stdErr, attachment)
	log.Println("Function call executed")

	// Post response to NATS
	functionReturn := functions.FunctionReturn{
		CallID: functionCall.CallID,
		Output: stdOutBuffer.Bytes(),
	}

	// Encode the proto object
	encoded, err := proto.Marshal(&functionReturn)
	if err != nil {
		log.Println("Error marshalling function return:", err)
		return
	}

	// Publish the return to NATS
	_, err = natsStreaming.GetClient().PublishAsync(functionCall.ReturnTo, encoded, nil)
	if err != nil {
		log.Println("Error publishing function return:", err)
		return
	}
	log.Println("Function return successfuly submitted, call Id:", functionReturn.CallID)
}

func containerAttach(ctx context.Context, containerID string) (types.HijackedResponse, error) {
	attachOptions := types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stream: true,
	}
	return dock.GetClient().ContainerAttach(ctx, containerID, attachOptions)
}

func containerStart(ctx context.Context, containerID string) error {
	startOptions := types.ContainerStartOptions{}
	return dock.GetClient().ContainerStart(ctx, containerID, startOptions)
}

func handleStreams(ctx context.Context, inputStream io.Reader, outputStream, errorStream *bufio.Writer, attachment types.HijackedResponse) error {
	receiveStdout := make(chan error, 1)
	if outputStream != nil || errorStream != nil {
		go func() {
			written, err := stdcopy.StdCopy(outputStream, errorStream, attachment.Reader)
			log.Printf("Transferred standard output (%v bytes)\n", written)
			receiveStdout <- err
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		written, err := io.Copy(attachment.Conn, inputStream)
		if err != nil {
			log.Printf("Couldn't copy input stream: %s\n", err)
		}
		if err := attachment.CloseWrite(); err != nil {
			log.Printf("Couldn't send EOF: %s\n", err)
		}
		log.Printf("Transferred standard input (%v bytes)\n", written)
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		if err != nil {
			log.Printf("Error receiveStdout: %s\n", err)
			return err
		}
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			select {
			case err := <-receiveStdout:
				if err != nil {
					log.Printf("Error receiveStdout: %s\n", err)
					return err
				}
			case <-ctx.Done():
			}
		}
	case <-ctx.Done():
	}
	outputStream.Flush()
	errorStream.Flush()

	return nil
}
