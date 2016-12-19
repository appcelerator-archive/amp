package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Unable to read standard input:", err)
	}
	log.Println("Got some input to title!") // Logs on standard error
	fmt.Print(strings.Title(string(input))) // Writes response on standard output
}
