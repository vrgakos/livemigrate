package migrate

import (
	"time"
	"fmt"
	"github.com/vrgakos/livemigrate/tcpapp"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/vrgakos/livemigrate/node"
	"sync"
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
	fmt.Printf("-- Milestone: %s (%v)\n", name, time.Now().Sub(m.LastTime))

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
	}

	b, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, job := range m.bgJobs {
		wg.Add(1)
		go (func(j *node.BackGroundjob) {
			j.Stop()
			wg.Done()
		})(job)
	}
	wg.Wait()

	os.Stdout.Sync()

	return ioutil.WriteFile(m.Opts.MeasureFileName, b, 0664)
}

func (m *Measure) AddStat(node *node.Node, file string) error {
	ssh, err := node.GetSshClient()
	if err != nil {
		return err
	}
	job := ssh.NewBackGroundjob("/mnt/docker/asd --file " + file)

	m.bgJobs = append(m.bgJobs, job)

	return nil
}