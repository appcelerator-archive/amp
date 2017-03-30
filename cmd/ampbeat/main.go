package main

import (
	"os"

	"github.com/appcelerator/amp/cmd/ampbeat/beater"
	"github.com/elastic/beats/libbeat/beat"
)

func main() {
	err := beat.Run("ampbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
