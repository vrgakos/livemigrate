package main

import (
	"github.com/akhenakh/statgo"
	"github.com/tealeg/xlsx"
	"fmt"
	"time"
	"log"
	"os"
	"os/signal"
)

func main() {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var err error


	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	log.Println("Start")


	run := true

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			run = false
		}
	}()

	row = sheet.AddRow()
	row.AddCell().SetValue("Time (s:ms)")
	row.AddCell().SetValue("User")
	row.AddCell().SetValue("Kernel")
	//row.AddCell().SetValue("Idle")
	row.AddCell().SetValue("IOWait")


	start := time.Now().UnixNano()
	for run {
		row = sheet.AddRow()

		row.AddCell().SetString(fmt.Sprintf("%.2f", (time.Now().UnixNano() - start) / 1000000000))

		s := statgo.NewStat()

		row.AddCell().SetValue(s.CPUStats().User)
		row.AddCell().SetValue(s.CPUStats().Kernel)
		//row.AddCell().SetValue(s.CPUStats().Idle)
		row.AddCell().SetValue(s.CPUStats().IOWait)
		//row.AddCell().SetValue(times[0].String())

		time.Sleep(time.Millisecond * 100)
	}

	log.Println("Received an interrupt, saving")
	err = file.Save("asd.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
