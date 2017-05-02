package migrate

import (
	"time"
	"fmt"
	"github.com/vrgakos/livemigrate/tcpapp"
	"github.com/tealeg/xlsx"
	"os"
	"github.com/vrgakos/livemigrate/node"
	"sync"
	"log"
)

type Milestone struct {
	Name		string
	At		time.Time
}

type Measure struct {
	Opts		*DoOpts

	StartTime	time.Time
	StopTime	time.Time
	LastTime	time.Time

	Milestones	[]*Milestone

	client		*tcpapp.TcpClient
	ClientResults	[]*tcpapp.TcpClientResult

	bgJobs		[]*node.BackGroundjob
}

func NewMeasure(opts *DoOpts) *Measure {
	return &Measure{
		Opts:		opts,
		Milestones:	make([]*Milestone, 0),
		ClientResults:  make([]*tcpapp.TcpClientResult, 0),
	}
}

func (m *Measure) MilestoneStart() {
	m.Milestones = make([]*Milestone, 0)
	m.LastTime = time.Now()
}

func (m *Measure) AddMilestone(name string) {
	stone := &Milestone{
		Name:		name,
		At:		time.Now(),
	}
	diff := time.Now().Sub(m.StartTime)
	secs := float64(diff.Nanoseconds()) / float64(1000000000)
	fmt.Printf("-- Milestone: %s (%v at %.1f)\n", name, time.Now().Sub(m.LastTime), secs)

	m.Milestones = append(m.Milestones, stone)
	m.LastTime = time.Now()
}


func (m *Measure) SetupClient(address string, interval time.Duration) error {
	m.ClientResults = make([]*tcpapp.TcpClientResult, 0)
	m.client = tcpapp.NewTcpClient(address, interval, func(res *tcpapp.TcpClientResult) {
		m.ClientResults = append(m.ClientResults, res)
	})
	return m.client.Start()
}

func (m *Measure) Start() error {
	m.Opts.Print()
	m.StartTime = time.Now()
	m.LastTime = time.Now()

	if len(m.Opts.TcpClientAddress) > 0 {
		err := m.SetupClient(m.Opts.TcpClientAddress, m.Opts.TcpClientInterval)
		if err != nil {
			return err
		}
	}

	time.Sleep(m.Opts.MeasureWaitBefore)
	return nil
}

func (m *Measure) Stop() error {
	time.Sleep(m.Opts.MeasureWaitAfter)

	m.StopTime = time.Now()
	if m.client != nil {
		m.client.Stop()

		/*b, err := json.MarshalIndent(m, "", "    ")
		if err != nil {
			return err
		}

		ioutil.WriteFile(m.Opts.MeasureFileName, b, 0664)*/

		xlsxFile := xlsx.NewFile()
		sheet, _ := xlsxFile.AddSheet("TCP RTT")
		row := sheet.AddRow()
		row.AddCell().SetValue("Time (s:ms)")
		row.AddCell().SetValue("RTT (ms)")

		for _, res := range m.ClientResults {
			timeString := fmt.Sprintf("%.1f", float64(res.Time.UnixNano() - m.StartTime.UnixNano()) / float64(1000000000))
			row := sheet.AddRow()
			row.AddCell().SetString(timeString)
			row.AddCell().SetFloat(float64(res.Rtt.Nanoseconds()) / float64(1000000))
		}
		xlsxFile.Save(m.Opts.MeasureFileName)
	}



	var wg sync.WaitGroup
	for _, job := range m.bgJobs {
		wg.Add(1)
		go (func(j *node.BackGroundjob) {
			j.Stop()
			//log.Println(b)
			//log.Println(err)
			wg.Done()
		})(job)
	}
	wg.Wait()

	m.AddMilestone("All done.")
	os.Stdout.Sync()

	return nil
}

func (m *Measure) AddStat(node *node.Node, file string) error {
	ssh, err := node.GetSshClient()
	if err != nil {
		log.Println(err)
		return err
	}
	job := ssh.NewBackGroundjob("/mnt/docker/asd --file " + file + " >> /tmp/asd.log 2>&1")

	m.bgJobs = append(m.bgJobs, job)

	return nil
}