package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/storage"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// StorageCmd is the main command for attaching storage subcommands.
var StorageCmd = &cobra.Command{
	Use:   "kv",
	Short: "Storage operations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return AMP.Connect()
	},
}

var (
	// storagePutCmd represents the creation of storage key-value pair
	storagePutCmd = &cobra.Command{
		Use:     "put",
		Short:   "Assign specified value with specified key",
		Example: "foo bar",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storagePut(AMP, args)
		},
	}
	// storageGetCmd represents the retrieval of storage value based on key
	storageGetCmd = &cobra.Command{
		Use:     "get",
		Short:   "Retrieve a storage object",
		Example: "foo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storageGet(AMP, args)
		},
	}
	// storageDeleteCmd represents the deletion of storage value based on key
	storageDeleteCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove a storage object",
		Example: "foo",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return storageDelete(AMP, args)
		},
	}
	// storageListCmd represents the list of storage key-value pair
	storageListCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List all storage objects",
		Example: "-q",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storageList(AMP, args)
		},
	}
)

func init() {
	RootCmd.AddCommand(StorageCmd)
	StorageCmd.AddCommand(storagePutCmd)
	StorageCmd.AddCommand(storageGetCmd)
	StorageCmd.AddCommand(storageDeleteCmd)
	StorageCmd.AddCommand(storageListCmd)
}

// storagePut validates the input command line arguments and creates or updates storage key-value pair
// by invoking the corresponding rpc/storage method
func storagePut(amp *cli.AMP, args []string) error {
	switch len(args) {
	case 0:
		mgr.Fatal("must specify storage key and storage value")
	case 1:
		mgr.Fatal("must specify storage value")
	case 2:
		// OK
	default:
		mgr.Fatal("too many arguments")
	}

	k := args[0]
	v := args[1]
	request := &storage.PutStorage{Key: k, Val: v}
	client := storage.NewStorageClient(amp.Conn)
	reply, err := client.Put(context.Background(), request)
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	fmt.Println(reply.Val)
	return nil
}

// storageGet validates the input command line arguments and retrieves storage key-value pair
//by invoking the corresponding rpc/storage method
func storageGet(amp *cli.AMP, args []string) error {
	if len(args) > 1 {
		mgr.Fatal("too many arguments")
	} else if len(args) == 0 {
		mgr.Fatal("must specify storage key")
	}
	k := args[0]
	if k == "" {
		mgr.Fatal("must specify storage key")
	}

	request := &storage.GetStorage{Key: k}

	client := storage.NewStorageClient(amp.Conn)
	reply, err := client.Get(context.Background(), request)
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	fmt.Println(reply.Val)
	return nil
}

// storageDelete validates the input command line arguments and deletes storage key-value pair
// by invoking the corresponding rpc/storage method
func storageDelete(amp *cli.AMP, args []string) error {
	if len(args) > 1 {
		mgr.Fatal("too many arguments")
	} else if len(args) == 0 {
		mgr.Fatal("must specify storage key")
	}
	k := args[0]
	if k == "" {
		mgr.Fatal("must specify storage key")
	}

	request := &storage.DeleteStorage{Key: k}

	client := storage.NewStorageClient(amp.Conn)
	reply, err := client.Delete(context.Background(), request)
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	fmt.Println(reply.Val)
	return nil
}

// storageList validates the input command line arguments and lists all the storage
// key-value pairs by invoking the corresponding rpc/storage method
func storageList(amp *cli.AMP, args []string) error {
	if len(args) > 0 {
		mgr.Fatal("too many arguments")
	}
	request := &storage.ListStorage{}
	client := storage.NewStorageClient(amp.Conn)
	reply, err := client.List(context.Background(), request)
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}
	if reply == nil || len(reply.List) == 0 {
		mgr.Warn("no storage object is available")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, tablePadding, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE\t")
	for _, info := range reply.List {
		fmt.Fprintf(w, "%s\t%s\t\n", info.Key, info.Val)
	}
	w.Flush()
	return nil
}
