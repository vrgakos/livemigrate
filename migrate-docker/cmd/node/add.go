package node

import (
	"github.com/spf13/cobra"
	"github.com/vrgakos/livemigrate/node"
	"fmt"
)

type addOptions struct {
	alias		string
	host		string
	sshUser 	string
	sshKeyFile	string

	sshPort		int
	dockerApiPort	int
}

func addCommand(store *node.NodeStore) *cobra.Command {
	var opts addOptions

	cmd := &cobra.Command{
		Use:   "add [OPTIONS] ALIAS HOST SSHUSER SSHKEYFILE",
		Short: "Add a node to config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return fmt.Errorf("Too few arguments")
			}

			opts.alias = args[0]
			opts.host = args[1]
			opts.sshUser = args[2]
			opts.sshKeyFile = args[3]

			node := store.NewNode(opts.alias, opts.host, opts.sshUser, opts.sshKeyFile)
			node.SshPort = opts.sshPort
			node.DockerApiPort = opts.dockerApiPort

			return store.Save()
		},
	}

	flags := cmd.Flags()
	flags.IntVar(&opts.sshPort, "ssh-port", 22, "Custom SSH tcp port")
	flags.IntVar(&opts.dockerApiPort, "docker-api-port", 2376, "Custom Docker API tcp port")

	return cmd
}
