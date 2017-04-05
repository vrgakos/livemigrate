package measure

import (
	"time"
	"fmt"
)

type Milestone struct {
	name		string
	at		time.Time
}

type Measure struct {
	start		time.Time
	last		time.Time

	milestones	[]*Milestone
}

func NewMeasure() *Measure {
	return &Measure{
		start:		time.Now(),
		last:		time.Now(),
		milestones:	make([]*Milestone, 0),
	}
}

func (m *Measure) MilestoneStart() {
	m.last = time.Now()
}

func (m *Measure) AddMilestone(name string) {
	stone := &Milestone{
		name:		name,
		at:		time.Now(),
	}
	fmt.Printf("-- Milestone: %s (%d ms)\n", name, time.Now().Sub(m.last).Nanoseconds() / 1000000)

	m.milestones = append(m.milestones, stone)
	m.last = time.Now()
}