package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Entry struct {
	Key   string `json:"key"`
	Flags int    `json:"flags"`
	Value string `json:"value"`
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s [encode/decode] FILE", os.Args[0])
	}

	// Parse JSON file
	raw, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	entries := make([]*Entry, 0)
	if err := json.Unmarshal(raw, &entries); err != nil {
		log.Fatal(err)
	}

	// Process JSON data
	switch os.Args[1] {
	case "encode":
		for _, entry := range entries {
			entry.Value = base64.StdEncoding.EncodeToString([]byte(entry.Value))
		}
	case "decode":
		for _, entry := range entries {
			decoded, _ := base64.StdEncoding.DecodeString(entry.Value)
			entry.Value = string(decoded)
		}
	default:
		log.Fatalf("Usage: %s [encode/decode] FILE", os.Args[0])
	}

	// Print output
	output, err := json.MarshalIndent(entries, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}
