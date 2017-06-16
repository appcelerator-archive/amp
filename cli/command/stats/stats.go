package stats

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type statsOpts struct {
	container       bool
	service         bool
	node            bool
	stack           bool
	cpu             bool
	mem             bool
	io              bool
	net             bool
	period          string
	timeGroup       string
	containerID     string
	containerName   string
	containerState  string
	serviceID       string
	stackName       string
	nodeID          string
	follow          bool
	includeAmpStats bool
}

var (
	opts = &statsOpts{}
)

var displayGroupMap = map[string]string{
	"container_short_name": "container",
	"service_name":         "service",
	"node_id":              "node",
	"stack_name":           "stack",
}

// NewStatsCommand returns a new instance of the stats command.
func NewStatsCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats [OPTIONS] SERVICE",
		Short: "Display amp statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getStats(c, args)
		},
	}
	flags := cmd.Flags()
	flags.BoolVar(&opts.container, "container", false, "Display stats on containers")
	flags.BoolVar(&opts.service, "service", true, "Display stats on services")
	flags.BoolVar(&opts.node, "node", false, "Display stats on nodes")
	flags.BoolVar(&opts.stack, "stack", false, "Display stats on stacks")
	//metrics
	flags.BoolVar(&opts.cpu, "cpu", false, "Display cpu stats")
	flags.BoolVar(&opts.mem, "mem", false, "Display memory stats")
	flags.BoolVar(&opts.io, "io", false, "Display disk io stats")
	flags.BoolVar(&opts.net, "net", false, "Display net rx/tx stats")
	//historic
	flags.StringVar(&opts.period, "period", "now-10m", `Historic period of metrics extraction, for instance: "now-1d", "now-10h", with y=year, M=month, w=week, d=day, h=hour, m=minute, s=second`)
	flags.StringVar(&opts.timeGroup, "time-group", "", `Historic extraction by time group, for instance: "1d", "3h", , with y=year, M=month, w=week, d=day, h=hour, m=minute, s=second`)
	//filters
	flags.StringVar(&opts.containerID, "container-id", "", "Filter on container id")
	flags.StringVar(&opts.containerName, "container-name", "", "Filter on container name")
	flags.StringVar(&opts.containerState, "container-state", "", "Filter on container state")
	flags.StringVar(&opts.serviceID, "service-id", "", "Filter on service id")
	flags.StringVar(&opts.stackName, "stack-name", "", "Filter on stack name")
	flags.StringVar(&opts.nodeID, "node-id", "", "Filter on node id")
	flags.BoolVarP(&opts.includeAmpStats, "include", "i", false, "Include AMP stats")
	//Stream flag
	flags.BoolVarP(&opts.follow, "follow", "f", false, "Follow stats output")
	return cmd
}

// Stats displays resource usage statistics
func getStats(c cli.Interface, args []string) error {
	var query = &stats.StatsRequest{}

	//Set filters
	query.FilterContainerId = opts.containerID
	query.FilterContainerName = opts.containerName
	query.FilterContainerState = opts.containerState
	query.FilterServiceId = opts.serviceID
	query.FilterStackName = opts.stackName
	query.FilterNodeId = opts.nodeID
	query.AllowsInfra = opts.includeAmpStats

	//set main Group
	if opts.timeGroup == "" {
		if len(args) > 0 {
			query.Group = "service_name"
			query.FilterServiceName = args[0]
		} else {
			if opts.container {
				query.Group = "container_short_name"
			} else if opts.node {
				query.Group = "node_id"
			} else if opts.stack {
				query.Group = "stack_name"
			} else {
				query.Group = "service_name"
			}
		}
	}

	//set metrics
	query.StatsCpu = opts.cpu
	query.StatsMem = opts.mem
	query.StatsIo = opts.io
	query.StatsNet = opts.net
	if !opts.cpu && !opts.mem && !opts.io && !opts.net {
		query.StatsCpu = true
		query.StatsMem = true
		query.StatsIo = true
		query.StatsNet = true
	}

	//Set historic parameters
	query.Period = opts.period
	query.TimeGroup = opts.timeGroup

	//set default for period only for current statistics
	if query.TimeGroup == "" && query.Period == "" {
		query.Period = "now-10s"
	}

	// Execute query regarding discriminator
	ctx := context.Background()
	conn := c.ClientConn()
	client := stats.NewStatsClient(conn)
	r, err := client.StatsQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if r.Entries == nil {
		fmt.Println("No result found")
		return nil
	}
	if query.TimeGroup != "" && query.Group != "" {
		fmt.Printf("%+v\n", r)
		return nil
	}
	if query.TimeGroup != "" {
		displayHistoricResult(query, r, c, true)
	} else {
		displayCurrentResult(query, r, c)
	}
	if opts.follow {
		if query.TimeGroup == "" {
			if err := continueCurrentExecution(ctx, c, client, query); err != nil {
				return err
			}
		} else {
			if err := continueHistoricExecution(ctx, c, client, query, r); err != nil {
				return err
			}
		}
	}
	return nil
}

// Display current stats
func displayCurrentResult(query *stats.StatsRequest, result *stats.StatsReply, c cli.Interface) {
	if opts.follow {
		c.Console().Printf("\033[2J\033[0;0H")
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	c.Console().Println(getGroupByFilterText(query))
	fmt.Fprintln(w, getStatTitle(query))
	for _, entry := range result.Entries {
		fmt.Fprintln(w, getOneStatLine(query, entry))
	}
	w.Flush()
}

// Display historic stats
func displayHistoricResult(query *stats.StatsRequest, result *stats.StatsReply, c cli.Interface, displayTitle bool) {
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	if displayTitle {
		c.Console().Println(getGroupByFilterText(query))
		fmt.Fprintln(w, getStatTitle(query))
	}
	for _, entry := range result.Entries {
		fmt.Fprintln(w, getOneStatLine(query, entry))
	}
	w.Flush()
}

// get one line of the stat table
func getOneStatLine(query *stats.StatsRequest, entry *stats.MetricsEntry) string {
	line := entry.Group
	if query.StatsCpu {
		line = fmt.Sprintf("%s\t%.2f%%", line, entry.Cpu.TotalUsage)
	}
	if query.StatsMem {
		line = fmt.Sprintf("%s\t%s\t%s\t%.1f%%", line, formatBytes(entry.Mem.Usage), formatBytes(entry.Mem.Limit), entry.Mem.UsageP*100)
	}
	if query.StatsIo {
		line = fmt.Sprintf("%s\t%s/s\t%s/s", line, formatBytes(entry.Io.Read/60), formatBytes(entry.Io.Write/60))
	}
	if query.StatsNet {
		line = fmt.Sprintf("%s\t%s/s\t%s/s", line, formatBytes(entry.Net.RxBytes/60), formatBytes(entry.Net.TxBytes/60))
	}
	return line
}

// get The title of the table to display
func getStatTitle(query *stats.StatsRequest) string {
	title := strings.ToUpper(getDisplayedGroup(query.Group))
	if query.TimeGroup != "" {
		title = "Date"
	}
	if query.StatsCpu {
		title = fmt.Sprintf("%s\tCPU %%", title)
	}
	if query.StatsMem {
		title = fmt.Sprintf("%s\tMEM USAGE\tLIMIT\tMEM %%", title)
	}
	if query.StatsIo {
		title = fmt.Sprintf("%s\tIO READ\tIO WRITE", title)
	}
	if query.StatsNet {
		title = fmt.Sprintf("%s\tNET RX\t NET TX", title)
	}
	return title
}

func getDisplayedGroup(group string) string {
	val, ok := displayGroupMap[group]
	if !ok {
		return "unknown"
	}
	return val
}

// get request explaination text to display at fist
func getGroupByFilterText(query *stats.StatsRequest) string {
	text := fmt.Sprintf("Stats on %ss period=%s", getDisplayedGroup(query.Group), query.Period)
	if query.TimeGroup != "" {
		text = fmt.Sprintf("Stats historic period=%s timeGroup=%s", query.Period, query.TimeGroup)
	}
	filters := []string{}
	if query.FilterContainerId != "" {
		filters = append(filters, fmt.Sprintf("ContainerId=%s", query.FilterContainerId))
	}
	if query.FilterContainerName != "" {
		filters = append(filters, fmt.Sprintf("ContainerName=%s", query.FilterContainerName))
	}
	if query.FilterContainerState != "" {
		filters = append(filters, fmt.Sprintf("ContainerState=%s", query.FilterContainerState))
	}
	if query.FilterServiceName != "" {
		filters = append(filters, fmt.Sprintf("ServiceName=%s", query.FilterServiceName))
	}
	if query.FilterServiceId != "" {
		filters = append(filters, fmt.Sprintf("ServiceId=%s", query.FilterServiceId))
	}
	if query.FilterStackName != "" {
		filters = append(filters, fmt.Sprintf("StackName=%s", query.FilterStackName))
	}
	if query.FilterTaskId != "" {
		filters = append(filters, fmt.Sprintf("TaskId=%s", query.FilterTaskId))
	}
	if query.FilterNodeId != "" {
		filters = append(filters, fmt.Sprintf("NodeId=%s", query.FilterNodeId))
	}
	if len(filters) == 0 {
		return fmt.Sprintf("%s, No filter", text)
	}
	for i, filter := range filters {
		if i == 0 {
			if len(filters) > 1 {
				text = fmt.Sprintf("%s, Filters: %s", text, filter)
			} else {
				text = fmt.Sprintf("%s, Filter: %s", text, filter)
			}
		} else {
			text = fmt.Sprintf("%s and %s", text, filter)
		}
	}
	return text
}

// format the number of bytes into number of Kb, Mb, Gb, ...
func formatBytes(vali int64) string {
	val := float64(vali)
	if val == 0 {
		return "0"
	} else if val < 1 {
		return "0.0"
	} else if val < 1024 {
		return fmt.Sprintf("%.1f B", val)
	} else if val < 1048576 {
		return fmt.Sprintf("%.1f KB", val/1024)
	} else if val < 1073741824 {
		return fmt.Sprintf("%.1f MB", val/1048576)
	}
	return fmt.Sprintf("%.1f GB", val/1073741824)
}

func continueCurrentExecution(ctx context.Context, c cli.Interface, client stats.StatsClient, query *stats.StatsRequest) error {
	for {
		time.Sleep(2 * time.Second)
		r, err := client.StatsQuery(ctx, query)
		if err != nil {
			return err
		}
		displayCurrentResult(query, r, c)
	}
}

func continueHistoricExecution(ctx context.Context, c cli.Interface, client stats.StatsClient, query *stats.StatsRequest, result *stats.StatsReply) error {
	sleepTime, err := getHistoricSleepDuration(query.TimeGroup)
	if err != nil {
		return err
	}
	lastDate := result.Entries[len(result.Entries)-1].Group
	currentDate, err := time.Parse("2006-01-02T15:04:05", lastDate)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", err)
	}
	for {
		time.Sleep(sleepTime)
		r, err := client.StatsQuery(ctx, query)
		if err != nil {
			return err
		}
		currentDate = currentDate.Add(sleepTime)
		if len(r.Entries) > 1 {
			r.Entries = r.Entries[len(r.Entries)-1 : len(r.Entries)]
			r.Entries[0].Group = currentDate.Format("2006-01-02T15:04:05")
		}
		displayHistoricResult(query, r, c, false)
	}
}

func getHistoricSleepDuration(timeGroup string) (time.Duration, error) {
	last := timeGroup[len(timeGroup)-1:]
	num, err := strconv.Atoi(timeGroup[0 : len(timeGroup)-1])
	if err != nil {
		return 0, err
	}
	if last == "y" || last == "M" || last == "w" || last == "d" {
		return 0, errors.New("cannot follow output when period is greater than a day")
	}
	if last == "h" {
		return time.Hour * time.Duration(num), nil
	} else if last == "m" {
		return time.Minute * time.Duration(num), nil
	}
	return time.Second * time.Duration(num), nil
}
