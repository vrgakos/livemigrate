package do

import (
	"github.com/spf13/cobra"
	"github.com/vrgakos/livemigrate/node"
	"github.com/vrgakos/livemigrate/migrate-docker/migrate"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Do migrate",
	/*Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("do called")
	},*/
}

func Init(root *cobra.Command, store *node.NodeStore) {
	opts := &migrate.DoOpts{}

	doCmd.PersistentFlags().IntVar(&opts.PredumpMaxIters, "max-iters", 8, "Number of maximum pre-dump operations")
	doCmd.PersistentFlags().IntVar(&opts.PredumpMinPages, "min-pages", 64, "Minimum memory pages per pre-dump")
	doCmd.PersistentFlags().IntVar(&opts.PredumpMaxGrowRate, "max-grow", 10, "Maximum allowed grow rate (percent) per pre-dump")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	doCmd.AddCommand(
		nfsCommand(store, opts),
	)
	root.AddCommand(doCmd)
}

