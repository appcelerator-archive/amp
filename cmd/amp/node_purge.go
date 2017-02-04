package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/spf13/cobra"
)

// NodePurgeCmd is the node purge command
var NodePurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "node purge",
	Long:  `node purge container images volumes`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nodePurge(AMP, cmd, args)
	},
}

func init() {
	NodeCmd.AddCommand(NodePurgeCmd)
	NodePurgeCmd.Flags().Bool("volume", false, "purge volumes")
	NodePurgeCmd.Flags().Bool("container", false, "purge containers")
	NodePurgeCmd.Flags().Bool("image", false, "purge images")
	NodePurgeCmd.Flags().BoolP("force", "f", false, "force purge")
	NodePurgeCmd.Flags().StringP("node", "n", "", "specify the node onto apply the purge, default all")
}

func nodePurge(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if err := amp.ConnectAdmServer(); err != nil {
		manager.fatalf("ampadmin_server is not available, start stack ampadmin first: 'amp pf start ampadmin':\n")
	}
	req := &servergrpc.PurgeNodesRequest{
		Node:      cmd.Flag("node").Value.String(),
		Container: cmd.Flag("container").Value.String() == "true",
		Volume:    cmd.Flag("volume").Value.String() == "true",
		Image:     cmd.Flag("image").Value.String() == "true",
		Force:     cmd.Flag("force").Value.String() == "true",
	}
	if !req.Container && !req.Volume && !req.Image {
		manager.fatalf("Nothing to purge please specify --container or --volume or --image\n")
	}
	client := servergrpc.NewClusterServerServiceClient(amp.ConnAdmServer)
	ret, err := client.PurgeNodes(ctx, req)
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	lines := []string{}
	for _, res := range ret.Agents {
		lines = append(lines, fmt.Sprintf("%s\t%d\t%d\t%d\t%s", res.AgentId, res.NbContainers, res.NbVolumes, res.NbImages, res.Error))
	}
	manager.displayInOrder("NODE\tREMOVED\tREMOVED\tREMOVED\tPURGE", "HOSTNAME (ID)\tCONTAINERS\tVOLUMES\tIMAGES\tSTATUS", lines)
	return nil
}
