package do

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/vrgakos/livemigrate/node"
	"github.com/vrgakos/livemigrate/migrate-docker/migrate"
)


func centralnfsCommand(store *node.NodeStore, migrateOpts *migrate.DoOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "centralnfs [OPTIONS] SOURCE DESTINATION NFS CONTAINER",
		Short: "Do live migrate over centralnfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return fmt.Errorf("Too few arguments!")
			}

			sourceNode := store.GetNode(args[0])
			if sourceNode == nil {
				return fmt.Errorf("Invalid source node!")
			}

			destNode := store.GetNode(args[1])
			if destNode == nil {
				return fmt.Errorf("Invalid destination node!")
			}

			nfsNode := store.GetNode(args[2])
			if destNode == nil {
				return fmt.Errorf("Invalid nfs node!")
			}

			migrate.CentralNfs(sourceNode, destNode, nfsNode, args[3], migrateOpts)


			return nil
		},
	}

	//flags := cmd.Flags()

	return cmd
}

