package cli

import (
  "fmt"
  "github.com/appcelerator/amp/api/rpc/stat"
  "github.com/appcelerator/amp/api/client"
  "github.com/spf13/cobra"
  "google.golang.org/grpc"
)


//CPU Display CPU Stats
func CPU(amp *client.AMP, cmd *cobra.Command) error {
  ctx, err := amp.GetAuthorizedContext()
  if err != nil {
    return err
  }
  if amp.Verbose() {
    fmt.Println("Cpu")
    fmt.Printf("Ressource: %v\n", cmd.Flag("ressourceName").Value)
  }
  conn, err := grpc.Dial(client.ServerAddress, grpc.WithInsecure())
  if err != nil {
    return err
  }

  request := stat.CPURequest{}
  request.RessourceName = cmd.Flag("RessourceName").Value.String()

  c := stat.NewStatClient(conn)
  r, err := c.CPUQuery(ctx, &request)
  if err != nil {
    return err
  }
  for _, entry := range r.Entries {
    fmt.Printf("%s %s\t\t%d", entry.ID, entry.Name, (entry.UsageUser*100)/entry.UsageTotal)
    //TODO format and
    //add (entry.UsageKernel*100)/entry.UsageTotal
    //add (entry.UsageSystem*100)/entry.UsageTotal

  }
  conn.Close()
  return nil
}
