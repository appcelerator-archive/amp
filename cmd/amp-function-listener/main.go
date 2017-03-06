package main

import (
	"fmt"
	"github.com/appcelerator/amp/data/functions"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/pkg/stringid"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/go-nats-streaming"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ## `amp-function-listener`
// This service role is to:
// - listen to HTTP events:
//    - Parse the HTTP body (if any) and use it as an input for the function
//    - Publish function call to NATS "function call" topic
//    - Wait on a channel with a timeout of one minute
//
// - listen to NATS for function returns on the "returnTo" topic. There is one "returnTo" topic per `amp-function-listener` used by workers to submit function return.
//   - Store the function return in a map
//   - Unblock the HTTP handler
//   - Get the output of the function (if any) and write it as the HTTP body
//   - In case of timeout, return an error

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// natsStreaming is the nats streaming client
	natsStreaming ns.NatsStreaming

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
	listenAddr = ":80"
)

func main() {
	log.Printf("%s (version: %s, build: %s)\n", os.Args[0], Version, Build)

	// Storage
	log.Println("Connecting to etcd at", amp.EtcdDefaultEndpoint)
	store = etcd.New([]string{amp.EtcdDefaultEndpoint}, "amp")
	if err := store.Connect(amp.DefaultTimeout); err != nil {
		log.Fatalln("Unable to connect to etcd:", err)
	}
	log.Println("Connected to etcd at", strings.Join(store.Endpoints(), ","))

	fStore = functions.NewStore(store)

	// NATS Connect
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname:", err)
	}
	if natsStreaming.Connect(amp.NatsDefaultURL, amp.NatsClusterID, os.Args[0]+"-"+hostname, amp.DefaultTimeout) != nil {
		log.Fatalln(err)
	}

	// NATS, subscribe to returnTo topic
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
	ctx, cancel := context.WithTimeout(context.Background(), amp.DefaultTimeout)
	defer cancel()
	function, err := fStore.GetFunction(ctx, idOrName)
	if err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching function: %s", err.Error()))
		return
	}
	if function == nil { // We didn't find the function by id, try by name (by listing them all)
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
	_, err = natsStreaming.GetClient().PublishAsync(amp.NatsFunctionTopic, encoded, nil)
	if err != nil {
		httpError(w, http.StatusInternalServerError, fmt.Sprintf("error publishing function call: %v", err))
		return
	}
	log.Println("Function call successfuly submitted, call id:", functionCall.CallID)

	// Wait for a NATS response
	locks[callID] = make(chan *functions.FunctionReturn, 1) // Create the channel
	select {
	case functionReturn := <-locks[callID]: // Wait for the functionReturn on the channel
		if _, err := fmt.Fprint(w, string(functionReturn.Output)); err != nil {
			httpError(w, http.StatusInternalServerError, fmt.Sprintf("error publishing function call: %v", err))
			return
		}
	case <-time.After(amp.DefaultTimeout): // Wait for timed out
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
