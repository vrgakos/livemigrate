package node

import (
	"github.com/spf13/cobra"
	"github.com/vrgakos/livemigrate/node"
	"fmt"
)

func listCommand(store *node.NodeStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls [OPTIONS]",
		Short: "List nodes from config file",
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Printf("ALIAS\tHOST\t\tSSHUSER\tSSHPORT\tDOCKERAPIPORT\n")
			for _, node := range store.Nodes {
				fmt.Printf("%s\t%s\t%s\t%d\t%d\n", node.Alias, node.Host, node.SshUser, node.SshPort, node.DockerApiPort)
			}
			
			return nil
		},
	}

	return cmd
}
