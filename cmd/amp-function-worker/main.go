package main

import (
	"bufio"
	"bytes"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/pkg/mq"
	"github.com/appcelerator/amp/pkg/mq/nats-streaming"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"os/signal"
)

// ## `amp-function-worker`
// This service role is to:
// - listen to function calls on the "function call" MQ topic
// - create the corresponding function container
// - attach to the created container (for stream management)
// - pass the function call input parameter through standard input
// - start the container
// - read the standard output of the container and wait for it to close
// - post the response back to MQ on the "returnTo" topic specified in the call

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// MQ is the message queuer interface
	MQ mq.Interface

	// docker is the Docker client
	docker *client.Client
)

func main() {
	log.Printf("%s (version: %s, build: %s)\n", os.Args[0], Version, Build)

	// Connect to message queuer
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname:", err)
	}
	MQ = ns.New(amp.NatsDefaultURL, amp.NatsClusterID, os.Args[0]+"-"+hostname)
	if err := MQ.Connect(amp.DefaultTimeout); err != nil {
		log.Fatal(err)
	}

	// Subscribe to function topic
	log.Println("Subscribing to topic:", amp.FunctionCallsQueue)
	_, err = MQ.Subscribe(amp.FunctionCallsQueue, messageHandler, &function.FunctionCall{}, mq.DeliverAllAvailable())
	if err != nil {
		MQ.Close()
		log.Fatalln("Unable to subscribe to topic", err)
	}
	log.Println("Subscribed to topic:", amp.FunctionCallsQueue)

	// Docker
	log.Printf("Connecting to Docker API at %s version API: %s\n", amp.DockerDefaultURL, amp.DockerDefaultVersion)
	docker, err = client.NewClient(amp.DockerDefaultURL, amp.DockerDefaultVersion, nil, nil)
	if err != nil {
		log.Fatalln("Unable to connect to docker", err)
	}
	log.Println("Connected to Docker API at", amp.DockerDefaultURL)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			log.Println("\nReceived an interrupt, unsubscribing and closing connection...")
			MQ.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func messageHandler(msg proto.Message, err error) {
	go processMessage(msg, err)
}

func processMessage(msg proto.Message, err error) {
	if err != nil {
		log.Println("Error in message processing:", err)
		return
	}

	fc, ok := msg.(*function.FunctionCall)
	if !ok {
		log.Println("Error in type assertion")
		return
	}

	log.Println("Function call received, call id:", fc.CallID)

	ctx := context.Background()

	// Create container
	container, err := containerCreate(ctx, fc.Function.Image)
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
	stdIn := bufio.NewReader(bytes.NewReader(fc.Input))

	// Standard output
	var stdOutBuffer bytes.Buffer
	stdOut := bufio.NewWriter(&stdOutBuffer)

	// Standard error
	var stdErrorBuffer bytes.Buffer
	stdErr := bufio.NewWriter(&stdErrorBuffer)

	// Handle standard streams
	streamCtx, cancel := context.WithTimeout(ctx, amp.DefaultTimeout/2)
	defer cancel()
	handleStreams(streamCtx, stdIn, stdOut, stdErr, attachment)
	log.Println("Function call executed")

	// Post response to MQ
	functionReturn := function.FunctionReturn{
		CallID: fc.CallID,
		Output: stdOutBuffer.Bytes(),
	}

	// Publish the return to MQ
	_, err = MQ.PublishAsync(fc.ReturnTo, &functionReturn, nil)
	if err != nil {
		log.Println("Error publishing function return:", err)
		return
	}
	log.Println("Function return successfuly submitted, call Id:", functionReturn.CallID)
}

func containerCreate(ctx context.Context, image string) (container.ContainerCreateCreatedBody, error) {
	containerConfig := &container.Config{
		Image:        image,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		StdinOnce:    true,
	}
	hostConfig := &container.HostConfig{}
	networkingConfig := &network.NetworkingConfig{}
	return docker.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, stringid.GenerateNonCryptoID())
}

func containerAttach(ctx context.Context, containerID string) (types.HijackedResponse, error) {
	attachOptions := types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stream: true,
	}
	return docker.ContainerAttach(ctx, containerID, attachOptions)
}

func containerStart(ctx context.Context, containerID string) error {
	startOptions := types.ContainerStartOptions{}
	return docker.ContainerStart(ctx, containerID, startOptions)
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
