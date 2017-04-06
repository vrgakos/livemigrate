package do

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/vrgakos/livemigrate/node"
	"github.com/vrgakos/livemigrate/migrate-docker/migrate"
)

type nfsOptions struct {
	source		string
	dest		string
	container	string

	maxIters	int
}

func nfsCommand(store *node.NodeStore) *cobra.Command {
	var opts nfsOptions

	cmd := &cobra.Command{
		Use:   "nfs [OPTIONS] SOURCE DESTINATION CONTAINER",
		Short: "Do live migrate over nfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 {
				return fmt.Errorf("Too few arguments!")
			}

			opts.source = args[0]
			opts.dest = args[1]
			opts.container = args[2]

			sourceNode := store.GetNode(opts.source)
			if sourceNode == nil {
				return fmt.Errorf("Invalid source!")
			}

			destNode := store.GetNode(opts.dest)
			if destNode == nil {
				return fmt.Errorf("Invalid destination!")
			}

			migrate.Nfs(sourceNode, destNode, opts.container, opts.maxIters)

			return store.Save()
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.maxIters, "max-iters", 5, "Max pre-dump count")

	return cmd
}

