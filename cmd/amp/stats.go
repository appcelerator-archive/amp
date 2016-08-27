package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const (
	blank     = "                                                                     "
	separator = "---------------------------------------------------------------------"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display resource usage statistics",
	Long:  `get statistics on containers, services, nodes about cpu, memory, io, net.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Stats(AMP, cmd, args)
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
	statsCmd.Flags().Bool("mem", false, "display memory stats")
	statsCmd.Flags().Bool("io", false, "display disk io stats")
	statsCmd.Flags().Bool("net", false, "display net rx/tx stats")
	//historic
	statsCmd.Flags().String("period", "", "historic period of metrics extraction, duration + time-group as 1m, 10m, 4h, see time-group")
	statsCmd.Flags().String("since", "", "date defining when begin the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm")
	statsCmd.Flags().String("until", "", "date defining when stop the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm")
	statsCmd.Flags().String("time-group", "", "historic extraction group can be: s:seconds, m:minutes, h:hours, d:days, w:weeks")
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
	//Stream flag
	statsCmd.Flags().BoolP("follow", "f", false, "Follow stat output")

	RootCmd.AddCommand(statsCmd)
}

// Stats displays resource usage statistcs
func Stats(amp *client.AMP, cmd *cobra.Command, args []string) error {
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		return err
	}

	var query = stats.StatsRequest{}

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

	//set metrics
	if cmd.Flag("cpu").Value.String() == "true" {
		query.StatsCpu = true
	}
	if cmd.Flag("mem").Value.String() == "true" {
		query.StatsMem = true
	}
	if cmd.Flag("io").Value.String() == "true" {
		query.StatsIo = true
	}
	if cmd.Flag("net").Value.String() == "true" {
		query.StatsNet = true
	}
	if !query.StatsCpu && !query.StatsMem && !query.StatsIo && !query.StatsNet {
		query.StatsCpu = true
		query.StatsMem = true
		query.StatsIo = true
		query.StatsNet = true
	}

	query.StatsFollow = false
	if cmd.Flag("follow").Value.String() == "true" {
		query.StatsFollow = true
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

	if amp.Verbose() {
		displayStatsQueryParameters(&query)
	}

	if err = validateQuery(&query); err != nil {
		return err
	}

	// Execute query regarding discriminator
	c := stats.NewStatsClient(AMP.Conn)

	if !query.StatsFollow {
		_, err = executeStat(ctx, c, &query, true, 0)
		return err
	}
	return startFollow(ctx, c, &query)
}

func validateQuery(query *stats.StatsRequest) error {
	if query.Period != "" && (query.Since != "" || query.Until != "") {
		return errors.New("--period can't be used with --since or --until")
	}
	return nil
}

func executeStat(ctx context.Context, c stats.StatsClient, query *stats.StatsRequest, title bool, currentTime int64) (int64, error) {
	r, err := c.StatsQuery(ctx, query)
	if err != nil {
		return 0, err
	}
	//fmt.Println(r.Entries[0].Time)
	if currentTime != 0 && r.Entries[0].Time == currentTime {
		return currentTime, nil
	}
	if query.Discriminator == "container" {
		displayContainer(query, r, title)
	} else if query.Discriminator == "service" {
		displayService(query, r, title)
	} else if query.Discriminator == "task" {
		displayTask(query, r, title)
	} else {
		displayNode(query, r, title)
	}
	return r.Entries[0].Time, nil
}

func displayStatsQueryParameters(query *stats.StatsRequest) {
	fmt.Println("Stat:")
	fmt.Printf("cpu:%t mem:%t io:%t net:%t on %s/n", query.StatsCpu, query.StatsMem, query.StatsIo, query.StatsNet, query.Discriminator)
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

func isHistoricQuery(req *stats.StatsRequest) bool {
	if req.Period != "" || req.Since != "" || req.Until != "" || req.TimeGroup != "" {
		return true
	}
	return false
}

func displayContainer(query *stats.StatsRequest, result *stats.StatsReply, title bool) {
	histoTitle, histoSub := getHistoColTitle(query)
	if query.FilterServiceName != "" {
		if title {
			fmt.Println("Service: " + query.FilterServiceName)
			fmt.Println(histoTitle + col("Container name", 40) + getMetricsTitle(query, ""))
			fmt.Println(histoSub + col("-", 20) + col("-", 40) + getMetricsTitle(query, "-"))
		}
		for _, row := range result.Entries {
			fmt.Println(getHistoColVal(query, row) + col(row.ContainerName, 40) + getMetricsCol(query, row))
		}
	} else {
		if title {
			fmt.Println(histoTitle + col("Service name", 20) + col("Container name", 40) + getMetricsTitle(query, ""))
			fmt.Println(histoSub + col("-", 25) + col("-", 20) + col("-", 40) + getMetricsTitle(query, "-"))
		}
		for _, row := range result.Entries {
			fmt.Println(getHistoColVal(query, row) + col(row.ServiceName, 20) + col(row.ContainerName, 40) + getMetricsCol(query, row))
		}
	}
}

func displayService(query *stats.StatsRequest, result *stats.StatsReply, title bool) {
	if title {
		histoTitle, histoSub := getHistoColTitle(query)
		fmt.Println(histoTitle + col("Service name", 20) + getMetricsTitle(query, ""))
		fmt.Println(histoSub + col("-", 20) + getMetricsTitle(query, "-"))
	}
	for _, row := range result.Entries {
		fmt.Println(getHistoColVal(query, row) + col(row.ServiceName, 20) + getMetricsCol(query, row))
	}
}

func displayTask(query *stats.StatsRequest, result *stats.StatsReply, title bool) {
	histoTitle, histoSub := getHistoColTitle(query)
	if query.FilterServiceName != "" {
		if title {
			fmt.Println("Service: " + query.FilterServiceName)
			fmt.Println(histoTitle + col("Task name", 25) + col("Node id", 30) + getMetricsTitle(query, ""))
			fmt.Println(histoSub + col("-", 25) + col("-", 30) + getMetricsTitle(query, "-"))
		}
		for _, row := range result.Entries {
			fmt.Println(getHistoColVal(query, row) + col(row.TaskName, 25) + col(row.NodeId, 30) + getMetricsCol(query, row))
		}
	} else {
		if title {
			fmt.Println(histoTitle + col("Service name", 20) + col("Task name", 25) + col("Node id", 30) + getMetricsTitle(query, ""))
			fmt.Println(histoSub + col("-", 20) + col("-", 25) + col("-", 30) + getMetricsTitle(query, "-"))
		}
		for _, row := range result.Entries {
			fmt.Println(getHistoColVal(query, row) + col(row.ServiceName, 25) + col(row.TaskName, 25) + col(row.NodeId, 30) + getMetricsCol(query, row))
		}
	}
}

func displayNode(query *stats.StatsRequest, result *stats.StatsReply, title bool) {
	if title {
		histoTitle, histoSub := getHistoColTitle(query)
		fmt.Println(histoTitle + col("Node id", 30) + getMetricsTitle(query, ""))
		fmt.Println(histoSub + col("-", 30) + getMetricsTitle(query, "-"))
	}
	for _, row := range result.Entries {
		fmt.Println(getHistoColVal(query, row) + col(row.NodeId, 30) + getMetricsCol(query, row))
	}
}

func getMetricsCol(query *stats.StatsRequest, row *stats.StatsEntry) string {
	var ret string
	if query.StatsCpu {
		ret = colr(fmt.Sprintf("%.1f%%", row.Cpu), 8)
	}
	if query.StatsMem {
		ret += colr(formatBytes(row.MemUsage), 12) + colr(fmt.Sprintf("%.1f%%", row.Mem), 8)
	}
	if query.StatsIo {
		ret += colr(formatBytes(row.IoRead), 10) + " / " + col(formatBytes(row.IoWrite), 10)
	}
	if query.StatsNet {
		ret += colr(formatBytes(row.NetRxBytes), 10) + " / " + col(formatBytes(row.NetTxBytes), 10)
	}
	return ret
}

func getMetricsTitle(query *stats.StatsRequest, un string) string {
	var ret string
	if query.StatsCpu {
		if un != "-" {
			ret = colr("CPU %%", 8)
		} else {
			ret = col("-", 8)
		}
	}
	if query.StatsMem {
		if un != "-" {
			ret += colr("Mem usage", 12) + colr("Mem %%", 8)
		} else {
			ret += col("-", 12) + col("-", 8)
		}
	}
	if query.StatsIo {
		if un != "-" {
			ret += colm("Disk IO read/write", 23)
		} else {
			ret += col("-", 23)
		}
	}
	if query.StatsNet {
		if un != "-" {
			ret += colm("Net Rx/Tx", 23)
		} else {
			ret += col("-", 23)
		}
	}
	return ret
}

func formatBytes(val float64) string {
	if val == 0 {
		return "0"
	} else if val < 1 {
		return "0.0"
	} else if val < 1024 {
		return fmt.Sprintf("%.0f B", val)
	} else if val < 1048576 {
		return fmt.Sprintf("%.1f KB", val/1024)
	} else if val < 1073741824 {
		return fmt.Sprintf("%.1f MB", val/1048576)
	}
	return fmt.Sprintf("%.1f GB", val/1073741824)
}

// display value in the left of a col
func col(value string, size int) string {
	if value == "-" {
		return separator[0:size]
	}
	if len(value) > size {
		return value[0:size]
	}
	return value + blank[0:size-len(value)]
}

// display value in the right of a col
func colr(value string, size int) string {
	if len(value) > size {
		return value[0:size]
	}
	return blank[0:size-len(value)] + value
}

// display value in the middle of a col
func colm(value string, size int) string {
	if len(value) > size {
		return value[0:size]
	}
	space := size - len(value)
	rest := space % 2
	return blank[0:space/2+rest] + value + blank[0:space/2]
}

// display time col
func colTime(val int64, size int) string {
	tm := time.Unix(val, 0)
	value := tm.Format("2006-01-02 15:04:05")
	return col(value, size)
}

func getHistoColTitle(query *stats.StatsRequest) (string, string) {
	if !isHistoricQuery(query) {
		return "", ""
	}
	return col("Time", 25), col("-", 25)
}

func getHistoColVal(query *stats.StatsRequest, row *stats.StatsEntry) string {
	if !isHistoricQuery(query) {
		return ""
	}
	/*
		if query.StatsFollow {
			return colTime(time.Now().Unix(), 25)
		}
	*/
	return colTime(row.Time, 25)
}

func startFollow(ctx context.Context, c stats.StatsClient, query *stats.StatsRequest) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	isHisto := isHistoricQuery(query)
	var currentTime int64
	if isHisto {
		ctime, err := executeStat(ctx, c, query, true, 0)
		currentTime = ctime
		if err != nil {
			return err
		}
		query.Since = ""
		query.Until = ""
		query.Period = ""
		query.StatsFollow = false
		if query.TimeGroup == "" {
			query.TimeGroup = "1m"
		}
		time.Sleep(1 * time.Second)
	}
	for {
		if !isHisto {
			fmt.Println("\033[0;0H")
		}
		ctime, err := executeStat(ctx, c, query, !isHisto, currentTime)
		currentTime = ctime
		if err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}
}
