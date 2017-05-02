package do

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/vrgakos/livemigrate/node"
	"github.com/vrgakos/livemigrate/migrate-docker/migrate"
)


func nfsCommand(store *node.NodeStore, migrateOpts *migrate.DoOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nfs [OPTIONS] SOURCE DESTINATION CONTAINER",
		Short: "Do live migrate over nfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return fmt.Errorf("Too few arguments!")
			}

			sourceNode := store.GetNode(args[0])
			if sourceNode == nil {
				return fmt.Errorf("Invalid source!")
			}

			destNode := store.GetNode(args[1])
			if destNode == nil {
				return fmt.Errorf("Invalid destination!")
			}

			migrate.Nfs(sourceNode, destNode, args[2], migrateOpts)


			return nil
		},
	}

	//flags := cmd.Flags()

	return cmd
}

