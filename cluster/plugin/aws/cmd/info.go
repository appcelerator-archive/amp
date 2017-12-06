package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/appcelerator/amp/cluster/plugin/aws/plugin"
	"golang.org/x/net/context"
)

func Info() {
	ctx := context.Background()
	resp, err := plugin.AWS.InfoStack(ctx)
	if err != nil {
		if j, jerr := plugin.PluginOutputToJSON(nil, nil, err); jerr == nil {
			// print json error to stdout
			fmt.Println(j)
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}

	j, err := plugin.PluginOutputToJSON(nil, resp, nil)
	if err != nil {
		log.Fatal(err)
	}

	// print json result to stdout
	fmt.Println(j)
}
