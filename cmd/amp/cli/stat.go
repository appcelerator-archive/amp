package cli

import (
  "fmt"
  "github.com/appcelerator/amp/api/rpc/stat"
  "github.com/appcelerator/amp/api/client"
  "github.com/spf13/cobra"
  "golang.org/x/net/context"
)


//Stat Display amp stats
func Stat(amp *client.AMP, cmd *cobra.Command) error {
  ctx, err := amp.GetAuthorizedContext()
  if err != nil {
    return err
  }

  var query = stat.StatRequest{}
  //set discriminator
  if (cmd.Flag("container").Value.String() == "true") {
    query.Discriminator = "container"
  } else if (cmd.Flag("service").Value.String() == "true") { 
    query.Discriminator = "service"
  } else {
    query.Discriminator = "node"
  }
  //set metric
  if (cmd.Flag("mem").Value.String() == "true") {
    query.Metric = "mem"
  } else if (cmd.Flag("io").Value.String() == "true") {
    query.Metric = "io"
  } else if (cmd.Flag("net").Value.String() == "true") {
    query.Metric = "net"
  } else {
    query.Metric = "cpu"
  }
  //Set filters
  query.FilterDatacenter = cmd.Flag("datacenter").Value.String()
  query.FilterHost = cmd.Flag("host").Value.String()
  query.FilterContainerId = cmd.Flag("container-id").Value.String()
  query.FilterContainerName = cmd.Flag("container-name").Value.String()
  query.FilterContainerImage = cmd.Flag("image").Value.String()
  query.FilterServiceId = cmd.Flag("service-id").Value.String()
  query.FilterServiceName = cmd.Flag("service-name").Value.String()
  query.FilterNodeId = cmd.Flag("node").Value.String()
  //Set historic parameters
  query.Period = cmd.Flag("period").Value.String()
  query.Since = cmd.Flag("since").Value.String()
  query.Until = cmd.Flag("until").Value.String()
  query.TimeUnit = cmd.Flag("time-unit").Value.String()

  if amp.Verbose() {
    displayStatQueryParameters(query)
  }

  //Execute query regarding discriminator
  c := stat.NewStatClient(amp.Connect())
  defer amp.Disconnect()
  if (query.Metric == "cpu") {
    return cpuStat(ctx, c, query)
  } else if (query.Metric == "mem") {
    return memStat(ctx, c, query)
  } else if (query.Metric == "io") {
    return ioStat(ctx, c, query) 
  } else if (query.Metric == "net") {
    return netStat(ctx, c, query)
  } 
  return nil
}


func cpuStat(ctx context.Context, c stat.StatClient, query stat.StatRequest) error {
  r, err := c.CPUQuery(ctx, &query)
  if err != nil {
    return err
  }
  //TODO: format ouput
  fmt.Println(r)
  return nil
}

func memStat(ctx context.Context, c stat.StatClient, query stat.StatRequest) error {
  r, err := c.MemQuery(ctx, &query)
  if err != nil {
    return err
  }
  //TODO: format ouput
  fmt.Println(r)
  return nil
}

func ioStat(ctx context.Context, c stat.StatClient, query stat.StatRequest) error {
  r, err := c.IOQuery(ctx, &query)
  if err != nil {
    return err
  }
  //TODO: format ouput
  fmt.Println(r)
  return nil
}

func netStat(ctx context.Context, c stat.StatClient, query stat.StatRequest) error {
  r, err := c.NetQuery(ctx, &query)
  if err != nil {
    return err
  }
  //TODO: format ouput
  fmt.Println(r)
  return nil
}

func displayStatQueryParameters(query stat.StatRequest) {
  fmt.Println("Stat:")
  fmt.Printf("metric: %v on %s/n", query.Metric, query.Discriminator)
  fmt.Println("filters:")
  if (query.FilterDatacenter != "") {
    fmt.Printf("datacenter = %v/n", query.FilterDatacenter)
  }
  if (query.FilterHost != "") {
    fmt.Printf("host = %v/n", query.FilterHost)
  }
  if (query.FilterContainerId != "") {
    fmt.Printf("container id = %v/n", query.FilterContainerId)
  }
  if (query.FilterContainerName != "") {
    fmt.Printf("container name = %v/n", query.FilterContainerName)
  }
  if (query.FilterContainerImage != "") {
    fmt.Printf("container image name = %v/n", query.FilterContainerImage)
  }
  if (query.FilterServiceId != "") {
    fmt.Printf("service id = %v/n", query.FilterServiceId)
  }
  if (query.FilterServiceName != "") {
    fmt.Printf("service name = %v/n", query.FilterServiceName)
  }
  if (query.FilterNodeId != "") {
    fmt.Printf("node id = %v/n", query.FilterNodeId)
  }
}

