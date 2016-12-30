package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

var origin = "http://localhost/"
var url = "ws://ws.websocket.local.atomiq.io/"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			default:
				var msg = make([]byte, 512)
				_, err = ws.Read(msg)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%s", msg)
			}
		}
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			default:
				var msg = make([]byte, 512)
				_, err := reader.Read(msg)
				if err != nil {
					log.Fatal(err)
				}
				ws.Write(msg)
			}
		}
	}()

	select {}
}
