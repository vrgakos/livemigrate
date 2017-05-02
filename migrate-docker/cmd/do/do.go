package do

import (
	"github.com/spf13/cobra"
	"github.com/vrgakos/livemigrate/node"
	"github.com/vrgakos/livemigrate/migrate-docker/migrate"
	"time"
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

	// Migrate options
	doCmd.PersistentFlags().IntVar(&opts.PredumpMaxIters, "max-iters", 8, "Number of maximum pre-dump operations")
	doCmd.PersistentFlags().IntVar(&opts.PredumpMinPages, "min-pages", 64, "Minimum memory pages per pre-dump")
	doCmd.PersistentFlags().IntVar(&opts.PredumpMaxGrowRate, "max-grow", 10, "Maximum allowed grow rate (percent) per pre-dump")


	// Measure options
	doCmd.PersistentFlags().StringVarP(&opts.MeasureFileName, "file", "f", "tcp-measure.xlsx", "Measure file path")
	doCmd.PersistentFlags().DurationVar(&opts.MeasureWaitBefore, "wait-before", time.Second * 1, "Wait before start migration in nanosec")
	doCmd.PersistentFlags().DurationVar(&opts.MeasureWaitAfter, "wait-after", time.Second * 1, "Wait after migration done in nanosec")

	// TCP client options
	doCmd.PersistentFlags().DurationVarP(&opts.TcpClientInterval, "client-interval", "i", time.Millisecond * 100, "TCP client message interval")
	doCmd.PersistentFlags().StringVarP(&opts.TcpClientAddress, "client-address", "a", "", "TCP client connection address")


	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	doCmd.AddCommand(
		nfsCommand(store, opts),
	)
	root.AddCommand(doCmd)
}

