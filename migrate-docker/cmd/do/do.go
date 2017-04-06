package do

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vrgakos/livemigrate/node"
)

// nodeCmd represents the node command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Do migrate",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("do called")
	},
}

func Init(root *cobra.Command, store *node.NodeStore) {
	doCmd.AddCommand(
		nfsCommand(store),
	)

	root.AddCommand(doCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

