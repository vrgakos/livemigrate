package migrate

import (
	"time"
	"fmt"
)

type DoOpts struct {
	PredumpMaxIters		int
	PredumpMinPages		int
	PredumpMaxGrowRate	int

	MeasureFileName		string
	MeasureWaitBefore	time.Duration
	MeasureWaitAfter	time.Duration

	TcpClientAddress	string
	TcpClientInterval	time.Duration
}

func (o *DoOpts) Print() {
	fmt.Printf("Migrate options:\n")
	fmt.Printf("  PredumpMaxIters = %d\n", o.PredumpMaxIters)
	fmt.Printf("  PredumpMinPages = %d\n", o.PredumpMinPages)
	fmt.Printf("  PredumpMaxGrowRate = %d\n", o.PredumpMaxGrowRate)
	fmt.Printf("Measure options:\n")
	fmt.Printf("  MeasureFileName = %s\n", o.MeasureFileName)
	fmt.Printf("  MeasureWaitBefore = %v\n", o.MeasureWaitBefore)
	fmt.Printf("  MeasureWaitAfter = %v\n", o.MeasureWaitAfter)
	fmt.Printf("TCP client options:\n")
	fmt.Printf("  TcpClientAddress = %s\n", o.TcpClientAddress)
	fmt.Printf("  TcpClientInterval = %v\n", o.TcpClientInterval)
}