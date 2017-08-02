package main

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/spf13/cobra"
)

const longForm = "Jan 2, 2006 3:04pm (MST)"

var (
	DefaultAddress = "http://localhost:9090"
	clientAPI      = *v1.httpAPI{}
)

func initClient(cmd *cobra.Command, args []string) {
	cfg := api.Config{
		Address: DefaultAddress,
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	clientAPI = v1.NewAPI(client)
}

func query(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	t, err := time.Parse(longForm, args[1])
	if err != nil {
		log.Fatal(err)
	}
	resp, err := clientAPI.Query(ctx, args[0], t)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Value type: %s, String: %s", resp.Type(), resp.String())
}

func queryRange(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	start, err := time.Parse(longForm, args[1])
	if err != nil {
		log.Fatal(err)
	}
	end, err := time.Parse(longForm, args[2])
	if err != nil {
		log.Fatal(err)
	}
	step, err := time.ParseDuration(args[3])
	if err != nil {
		log.Fatal(err)
	}
	queryRange := &v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	}
	resp, err := clientAPI.QueryRange(ctx, args[0], queryRange)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Value type: %s, String: %s", resp.Type(), resp.String())
}

func labelValues(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	resp, err := clientAPI.LabelValues(ctx, args[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Value type: %s, String: %s", resp.Type(), resp.String())
}

func main() {
	rootCmd := &cobra.Command{
		Use:              "prometheus",
		Short:            "Query prometheus for cluster metrics",
		PersistentPreRun: initClient,
	}

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "evaluate an instant query at a given point in time",
		Run:   query,
	}

	queryRangeCmd := &cobra.Command{
		Use:   "queryrange",
		Short: "evaluates an expression query over a range of time",
		Run:   queryRange,
	}

	labelValuesCmd := &cobra.Command{
		Use:   "labelvalues",
		Short: "returns  the list of time series that match a certain label set",
		Run:   labelValues,
	}

	rootCmd.AddCommand(queryCmd, queryRangeCmd, labelValuesCmd)

	_ = rootCmd.Execute()
}
