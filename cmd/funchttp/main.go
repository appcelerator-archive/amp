package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/appcelerator/amp/data/functions"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
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

	// store is the interface used to access the key/value storage backend
	store storage.Interface

	// functions is the interface used to access the function storage
	fStore functions.Interface

	// returnToTopic is the topic used to listen to function return
	returnToTopic string

	// locks is used for function return (indexed by call id)
	locks = make(map[string](chan *functions.FunctionReturn))
)

const (
	name                = "funchttp"
	listenAddr          = ":80"
	connectionTimeout   = 5 * time.Second
	functionGetTimeout  = time.Minute
	functionCallTimeout = time.Minute
)

func main() {
	log.Printf("%s (version: %s, build: %s)\n", name, Version, Build)

	// Storage
	store = etcd.New([]string{etcd.DefaultEndpoint}, "amp", connectionTimeout)
	log.Println("Connecting to etcd at", etcd.DefaultEndpoint)
	if err := store.Connect(); err != nil {
		log.Fatalln("Unable to connect to etcd:", err)
	}
	log.Println("Connected to etcd at", strings.Join(store.Endpoints(), ","))

	fStore = functions.NewStore(store)

	// Nats
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname:", err)
	}
	natsStreaming = ns.NewClient(ns.DefaultURL, ns.ClusterID, name+"-"+hostname, connectionTimeout)
	if err := natsStreaming.Connect(); err != nil {
		log.Fatalln(err)
	}

	// Subscribe to returnTo topic
	returnToTopic = "returnTo-" + hostname
	log.Println("Subscribing to topic:", returnToTopic)
	_, err = natsStreaming.GetClient().Subscribe(returnToTopic, messageHandler, stan.DeliverAllAvailable())
	if err != nil {
		natsStreaming.Close()
		log.Fatalln("Unable to subscribe to topic", err)
	}
	log.Println("Subscribed to topic:", returnToTopic)

	// HTTP
	router := httprouter.New()
	router.POST("/:function", Index)

	log.Println("Start listening on", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, router))
}

// Index index
func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Get the function id or name from the URL
	idOrName := strings.TrimSpace(p.ByName("function"))

	// Try to get the function by id first
	ctx, cancel := context.WithTimeout(context.Background(), functionGetTimeout)
	defer cancel()
	function, err := fStore.GetFunction(ctx, idOrName)
	if err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching function: %s", err.Error()))
		return
	}
	if function == nil {
		// We didn't find the function by id, try by name (by listing them all)
		function, err = fStore.GetFunctionByName(ctx, idOrName)
		if err != nil {
			httpError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching function: %s", err.Error()))
			return
		}
		if function == nil {
			httpError(w, http.StatusNotFound, fmt.Sprintf("function not found: %s", idOrName))
			return
		}
	}
	log.Println("Function found", function)

	// Read the body parameter if any
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpError(w, http.StatusBadRequest, fmt.Sprintf("error reading request body: %v", err))
		return
	}

	// Invoke the function by posting a function call to NATS
	callID := stringid.GenerateNonCryptoID()
	functionCall := functions.FunctionCall{
		CallID:   callID,
		Input:    body,
		Function: function,
		ReturnTo: returnToTopic,
	}

	// Encode the proto object
	encoded, err := proto.Marshal(&functionCall)
	if err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Sprintf("error marshalling function call: %v", err))
		return
	}

	// Publish to NATS
	_, err = natsStreaming.GetClient().PublishAsync(ns.FunctionSubject, encoded, nil)
	if err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Sprintf("error publishing function call: %v", err))
		return
	}
	log.Println("Function call successfuly submitted, call id:", functionCall.CallID)

	// Wait for a NATS response
	locks[callID] = make(chan *functions.FunctionReturn, 1) // Create the channel
	select {
	case functionReturn := <-locks[callID]:
		// Wait for the functionReturn on the channel
		if _, err := fmt.Fprint(w, string(functionReturn.Output)); err != nil {
			httpError(w, http.StatusInternalServerError, fmt.Sprintf("error publishing function call: %v", err))
			return
		}
	case <-time.After(functionCallTimeout):
		// Wait for timed out
		httpError(w, http.StatusRequestTimeout, "function call timed out")
		return
	}
}

func httpError(w http.ResponseWriter, statusCode int, message string) {
	log.Println(message)
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, message)
}

func messageHandler(msg *stan.Msg) {
	// Parse function return message
	functionReturn := &functions.FunctionReturn{}
	err := proto.Unmarshal(msg.Data, functionReturn)
	if err != nil {
		log.Println("Error unmarshalling function return:", err)
		return
	}
	log.Println("Function return received, call id:", functionReturn.CallID)

	// Unlock the caller
	locks[functionReturn.CallID] <- functionReturn
}
