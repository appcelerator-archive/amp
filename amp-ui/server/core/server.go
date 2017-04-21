package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"
)

//Server data
type Server struct {
	api *serverAPI
}

//ServerInit Connect to docker engine, get initial containers list and start the agent
func ServerInit(version string, build string) error {
	server := Server{api: &serverAPI{}}
	server.trapSignal()
	conf.init(version, build)
	server.start()
	return nil
}

// Launch a routine to catch SIGTERM Signal
func (s *Server) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\namp-ui server received SIGTERM signal")
		os.Exit(1)
	}()
}

func (s *Server) start() {
	r := mux.NewRouter()
	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(r)

	abspath, err := filepath.Abs("./public")
	if err != nil {
		fmt.Print(err)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(abspath)))
	s.api.handleAPIFunctions(r)
	log.Printf("AMP-UI server starting on %s\n", conf.port)
	if err := http.ListenAndServe(":"+conf.port, n); err != nil {
		log.Fatal("Server error: ", err)
	}
}
