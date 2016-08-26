package main

import (
	"fmt"
	"strconv"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stat"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const blank = "                                                                     "

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display resource usage statistics",
	Long:  `get statistics on containers, services, nodes about cpu, memory, io, net.`,
	Run: func(cmd *cobra.Command, args []string) {
		amp := client.NewAMP(&Config)
		err := Stats(amp, cmd, args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	statsCmd.Flags().Bool("container", false, "display stats on containers")
	statsCmd.Flags().Bool("service", false, "displat stats on services")
	statsCmd.Flags().Bool("node", false, "display stats on nodes")
	statsCmd.Flags().Bool("task", false, "display stats on tasks")
	//metrics
	statsCmd.Flags().Bool("cpu", false, "display cpu stats")
	//statsCmd.Flags().Bool("mem", false, "display memory stats")
	//statsCmd.Flags().Bool("io", false, "display memory stats")
	//statsCmd.Flags().Bool("net", false, "display memory stats")
	//historic
	statsCmd.Flags().String("period", "", "historic period of metrics extraction, duration + time unit")
	statsCmd.Flags().String("since", "", "date defining when begin the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm")
	statsCmd.Flags().String("until", "", "date defining when stop the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm")
	//statsCmd.Flags().String("time-unit", "", "historic extraction group can be: s:seconds, m:minutes, h:hours, d:days, w:weeks")
	//filters:
	statsCmd.Flags().String("container-id", "", "filter on container id")
	statsCmd.Flags().String("container-name", "", "filter on container name")
	statsCmd.Flags().String("image", "", "filter on container image name")
	statsCmd.Flags().String("service-name", "", "filter on service name")
	statsCmd.Flags().String("service-id", "", "filter on service id")
	statsCmd.Flags().String("task-name", "", "filter on task name")
	statsCmd.Flags().String("task-id", "", "filter on task id")
	statsCmd.Flags().String("datacenter", "", "filter on datacenter")
	statsCmd.Flags().String("host", "", "filter on host")
	statsCmd.Flags().String("node-id", "", "filter on node id")

	RootCmd.AddCommand(statsCmd)
}

// Stats displays resource usage statistcs
func Stats(amp *client.AMP, cmd *cobra.Command, args []string) error {
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		return err
	}

	var query = stat.StatRequest{}
	//set discriminator
	if cmd.Flag("container").Value.String() == "true" {
		query.Discriminator = "container"
	} else if cmd.Flag("service").Value.String() == "true" {
		query.Discriminator = "service"
	} else if cmd.Flag("task").Value.String() == "true" {
		query.Discriminator = "task"

	} else {
		query.Discriminator = "node"
	}

	//Set filters
	query.FilterDatacenter = cmd.Flag("datacenter").Value.String()
	query.FilterHost = cmd.Flag("host").Value.String()
	query.FilterContainerId = cmd.Flag("container-id").Value.String()
	query.FilterContainerName = cmd.Flag("container-name").Value.String()
	query.FilterContainerImage = cmd.Flag("image").Value.String()
	query.FilterServiceId = cmd.Flag("service-id").Value.String()
	query.FilterServiceName = cmd.Flag("service-name").Value.String()
	query.FilterTaskId = cmd.Flag("task-id").Value.String()
	query.FilterTaskName = cmd.Flag("task-name").Value.String()
	query.FilterNodeId = cmd.Flag("node-id").Value.String()
	//Set historic parameters
	query.Period = cmd.Flag("period").Value.String()
	query.Since = cmd.Flag("since").Value.String()
	query.Until = cmd.Flag("until").Value.String()
	//query.TimeUnit = cmd.Flag("time-unit").Value.String()

	if amp.Verbose() {
		displayStatQueryParameters(&query)
	}

	if err = validateQuery(&query); err != nil {
		return err
	}

	//Execute query regarding discriminator
	c := stat.NewStatClient(amp.Connect())
	defer amp.Disconnect()

	return cpuStat(ctx, c, &query)
}

func validateQuery(query *stat.StatRequest) error {
	// TODO consider implementing basic query validation here before calling service
	return nil
}

func cpuStat(ctx context.Context, c stat.StatClient, query *stat.StatRequest) error {
	r, err := c.CPUQuery(ctx, query)
	if err != nil {
		return err
	}
	if query.Discriminator == "container" {
		displayCPUContainer(query, r)
	} else if query.Discriminator == "service" {
		displayCPUService(query, r)
	} else if query.Discriminator == "task" {
		displayCPUService(query, r)
	} else if query.Discriminator == "node" {
		displayCPUService(query, r)
	} else {
		displayCPUNode(query, r)
	}
	return nil
}

func displayStatQueryParameters(query *stat.StatRequest) {
	fmt.Println("Stat:")
	fmt.Printf("metric: %v on %s/n", query.Metric, query.Discriminator)
	fmt.Println("filters:")
	if query.FilterDatacenter != "" {
		fmt.Printf("datacenter = %v/n", query.FilterDatacenter)
	}
	if query.FilterHost != "" {
		fmt.Printf("host = %v/n", query.FilterHost)
	}
	if query.FilterContainerId != "" {
		fmt.Printf("container id = %v/n", query.FilterContainerId)
	}
	if query.FilterContainerName != "" {
		fmt.Printf("container name = %v/n", query.FilterContainerName)
	}
	if query.FilterContainerImage != "" {
		fmt.Printf("container image name = %v/n", query.FilterContainerImage)
	}
	if query.FilterServiceId != "" {
		fmt.Printf("service id = %v/n", query.FilterServiceId)
	}
	if query.FilterServiceName != "" {
		fmt.Printf("service name = %v/n", query.FilterServiceName)
	}
	if query.FilterTaskId != "" {
		fmt.Printf("task id = %v/n", query.FilterTaskId)
	}
	if query.FilterTaskName != "" {
		fmt.Printf("task name = %v/n", query.FilterTaskName)
	}
	if query.FilterNodeId != "" {
		fmt.Printf("node id = %v/n", query.FilterNodeId)
	}
}

func displayCPUContainer(query *stat.StatRequest, result *stat.CPUReply) {
	for _, row := range result.Entries {
		fmt.Println(colTime(row.Time, 25) + col(row.ContainerName, 50) + col(row.NodeId, 30) + getCPUCols(row))
	}
}

func displayCPUService(query *stat.StatRequest, result *stat.CPUReply) {
	for _, row := range result.Entries {
		fmt.Println(colTime(row.Time, 25) + col(row.ServiceName, 20) + col(row.NodeId, 30) + getCPUCols(row))
	}
}

func displayCPUTask(query *stat.StatRequest, result *stat.CPUReply) {
	for _, row := range result.Entries {
		fmt.Println(colTime(row.Time, 25) + col(row.TaskName, 20) + col(row.NodeId, 30) + getCPUCols(row))
	}
}

func displayCPUNode(query *stat.StatRequest, result *stat.CPUReply) {
	for _, row := range result.Entries {
		fmt.Println(colTime(row.Time, 25) + col(row.Datacenter, 20) + col(row.Host, 30) + col(row.NodeId, 30) + getCPUCols(row))
	}
}

func getCPUCols(row *stat.CPUEntry) string {
	//usageSystem, _ := strconv.ParseFloat(row.UsageSystem, 64)
	//usageKernel, _ := strconv.ParseFloat(row.UsageKernel, 64)
	usageUser, _ := strconv.ParseFloat(row.UsageUser, 64)
	usageTotal, _ := strconv.ParseFloat(row.UsageTotal, 64)
	//var system string
	//var kernel string
	var user string
	if usageTotal != 0 {
		//system = fmt.Sprintf("%f", usageSystem * 100 / usageTotal)
		//kernel = fmt.Sprintf("%f", usageKernel * 100 / usageTotal)
		user = fmt.Sprintf("%.1f", usageUser*100/usageTotal)
	}
	//return col(system, 12) + col(kernel, 12) + col(user, 12)
	return col(user, 12)
}

func col(value string, size int) string {
	if len(value) > size {
		return value[0:size]
	}
	return value + blank[0:size-len(value)]
}

func colTime(val int64, size int) string {
	value := fmt.Sprintf("%d", val)
	if len(value) > size {
		return value[0:size]
	}
	return value + blank[0:size-len(value)]
}
