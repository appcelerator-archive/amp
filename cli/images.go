package cli

import (
	"fmt"
	flag "github.com/spf13/pflag"
)

var target string

func main() {

	parseFlags()

	fmt.Println("target environment:", target)

	images, err := LoadImageList()
	if err != nil {
		fmt.Println("unable to load config/images.yml (are you running from ampswarm root?)")
	}

	for _, v := range images {
		fmt.Println(v)
	}
}

func parseFlags() {
	flag.StringVarP(&target, "target", "t", "local", "target swarm environment (local|virtualbox|aws)")
}
