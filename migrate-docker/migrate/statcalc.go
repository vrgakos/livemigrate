package migrate

import (
	"github.com/docker/docker/api/types"
)

type StatCalc struct {
	prevStats   *types.CheckpointStat
	currStats   *types.CheckpointStat
	opts        *DoOpts

	iters	    int
}

func NewStatCalc(opts *DoOpts) *StatCalc {
	return &StatCalc{
		opts:	opts,
		iters:  0,
	}
}

func (s *StatCalc) Add(stats *types.CheckpointStat) {
	s.iters++
	s.prevStats = s.currStats
	s.currStats = stats

	//log.Printf("Pages written=%d, scanned=%d, skipped=%d\n", stats.PagesWritten, stats.PagesScanned, stats.PagesSkippedParent)
	//log.Printf("memdumpTime=%d, memwriteTime=%d\n", stats.MemdumpTime, stats.MemwriteTime)
	//log.Printf("frozenTime=%d, freezingTime=%d\n", stats.FrozenTime, stats.FreezingTime)
}

func (s *StatCalc) Resume() bool {
	if s.iters >= s.opts.PredumpMaxIters {
		return false
	}

	if s.currStats != nil {
		if s.currStats.PagesWritten <= uint64(s.opts.PredumpMinPages) {
			return false
		}

		if s.prevStats != nil {
			growRate := int64(s.currStats.PagesWritten - s.prevStats.PagesWritten) / int64(s.prevStats.PagesWritten / 100)
			if growRate > int64(s.opts.PredumpMaxGrowRate) {
				return false
			}
		}
	}

	return true
}

func (s *StatCalc) GetIters() int {
	return s.iters
}