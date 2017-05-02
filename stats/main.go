package main

import (
	"github.com/akhenakh/statgo"
	"github.com/tealeg/xlsx"
	"fmt"
	"time"
	"log"
	"os"
	"os/signal"
	"flag"
	"path/filepath"
)

var fileName string
var interval time.Duration

func init() {
	flag.StringVar(&fileName, "file", "stats.xlsx", "Output file path")
	flag.DurationVar(&interval, "interval", time.Millisecond * 100, "Measure interval")
}

func main() {
	flag.Parse()
	dir, _ := filepath.Abs(filepath.Dir(fileName))
	log.Println("Create directory", dir)
	os.MkdirAll(dir, 0777)

	var file *xlsx.File
	var sheetCpu, sheetMem *xlsx.Sheet
	var netSheets []*xlsx.Sheet
	var diskSheets []*xlsx.Sheet
	var row *xlsx.Row
	var err error


	file = xlsx.NewFile()
	sheetCpu, _ = file.AddSheet("CPU")
	sheetMem, _ = file.AddSheet("MEM")

	log.Println("Start")


	run := true

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			run = false
		}
	}()

	// CPU
	row = sheetCpu.AddRow()
	row.AddCell().SetValue("Time (s:ms)")
	row.AddCell().SetValue("User")
	row.AddCell().SetValue("Kernel")
	//row.AddCell().SetValue("Idle")
	row.AddCell().SetValue("IOWait")

	// MEM
	row = sheetMem.AddRow()
	row.AddCell().SetValue("Time (s:ms)")
	row.AddCell().SetValue("Used")
	row.AddCell().SetValue("Free")
	row.AddCell().SetValue("Cache")

	// NET
	s := statgo.NewStat()
	for _, stat := range s.NetIOStats() {
		sheet, _ := file.AddSheet("NET-" + stat.IntName)
		netSheets = append(netSheets, sheet)

		row = sheet.AddRow()
		row.AddCell().SetValue("Time (s:ms)")
		row.AddCell().SetValue("RX (MB/s)")
		row.AddCell().SetValue("TX (MB/s)")
		row.AddCell().SetValue("IPackets (packets/s)")
		row.AddCell().SetValue("OPackets (packets/s)")
		row.AddCell().SetValue("seconds")
	}

	// DISK
	for _, stat := range s.DiskIOStats() {
		sheet, _ := file.AddSheet("DISK-" + stat.DiskName)
		diskSheets = append(diskSheets, sheet)

		row = sheet.AddRow()
		row.AddCell().SetValue("Time (s:ms)")
		row.AddCell().SetValue("Read (MB/s)")
		row.AddCell().SetValue("Write (MB/s)")
		row.AddCell().SetValue("seconds")
	}




	start := time.Now().UnixNano()
	last := time.Now()
	for run {
		timer := time.NewTimer(interval)

		diff := time.Now().Sub(last)
		last = time.Now()
		secs := float64(diff.Nanoseconds()) / float64(1000000000)

		timeString := fmt.Sprintf("%.1f", float64(time.Now().UnixNano() - start) / float64(1000000000))


		s := statgo.NewStat()

		// CPU STATS
		cpuStats := s.CPUStats()
		row = sheetCpu.AddRow()
		row.AddCell().SetString(timeString)
		row.AddCell().SetFloat(cpuStats.User)
		row.AddCell().SetFloat(cpuStats.Kernel)
		//row.AddCell().SetValue(s.CPUStats().Idle)
		row.AddCell().SetFloat(cpuStats.IOWait)


		// MEM STATS
		memStats := s.MemStats()
		row = sheetMem.AddRow()
		row.AddCell().SetString(timeString)
		row.AddCell().SetFloat(float64(memStats.Used) / 1024 / 1024)
		row.AddCell().SetFloat(float64(memStats.Free) / 1024 / 1024)
		row.AddCell().SetFloat(float64(memStats.Cache) / 1024 / 1024)


		// NET STATS
		for i, stat := range s.NetIOStats() {
			sheet := netSheets[i]

			if stat.Period.Seconds() > 0 {
				row = sheet.AddRow()
				secs = stat.Period.Seconds()
				row.AddCell().SetString(timeString)
				row.AddCell().SetFloat((float64(stat.RX) / 1024 / 1024) / secs)
				row.AddCell().SetFloat((float64(stat.TX) / 1024 / 1024) / secs)
				row.AddCell().SetFloat(float64(stat.IPackets) / secs)
				row.AddCell().SetFloat(float64(stat.OPackets) / secs)
				row.AddCell().SetFloat(secs)
			}
		}

		// DISK STATS
		for i, stat := range s.DiskIOStats() {
			sheet := diskSheets[i]

			if stat.Period.Seconds() > 0 {
				row = sheet.AddRow()
				secs = stat.Period.Seconds()
				row.AddCell().SetString(timeString)
				row.AddCell().SetFloat(float64(stat.ReadBytes) / 1024 / 1024 / secs)
				row.AddCell().SetFloat(float64(stat.WriteBytes) / 1024 / 1024 / secs)
				row.AddCell().SetFloat(secs)
			}
		}

		<-timer.C
	}

	log.Println("Received an interrupt, saving")
	err = file.Save(fileName)
	if err != nil {
		fmt.Printf(err.Error())
	}
}
