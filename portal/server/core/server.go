package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"
)

//Server data
type Server struct {
	conf *ServerConfig
	api  *serverAPI
}

//ServerInit Connect to docker engine, get initial containers list and start the agent
func ServerInit(version string, build string) error {
	server := Server{api: &serverAPI{}}
	server.conf = &ServerConfig{}
	server.api.conf = server.conf
	server.trapSignal()
	server.conf.init(version, build)
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
	log.Println("Waiting for amplifier ready...")
	err := s.connectAmplifier()
	for err != nil {
		time.Sleep(5 * time.Second)
		err = s.connectAmplifier()
	}
	log.Printf("connected to amplifier: %s\n", s.conf.amplifierAddr)
	r := mux.NewRouter()
	n := negroni.Classic()
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(r)

	abspath, err := filepath.Abs("./public")
	if err != nil {
		fmt.Print(err)
	}
	s.api.handleAPIFunctions(r)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(abspath)))
	log.Printf("AMP-UI server starting on %s\n", s.conf.port)
	if err := http.ListenAndServe(":"+s.conf.port, n); err != nil {
		log.Fatal("Server error: ", err)
	}
}

func (s *Server) connectAmplifier() error {
	conn, err := grpc.Dial(s.conf.amplifierAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*3),
	)
	if err != nil {
		return err
	}
	s.api.conn = conn
	return nil
}
