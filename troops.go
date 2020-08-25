package troops

import (
	"time"
)

const (
	scaleMax = 3
)

type Troops struct {
	MaxSize        int
	MinIdle        int
	MaxIdleTimeout time.Duration
	JobBuffer      chan *Job
	Soldiers       chan IWorker
	WorkingSize    int
	Log            ILogger
}

func NewTroops(maxSize, minIdle int) *Troops {
	return &Troops{
		MaxSize:        maxSize,
		MinIdle:        minIdle,
		MaxIdleTimeout: 10 * time.Minute,
		JobBuffer:      make(chan *Job, maxSize*scaleMax),
		Soldiers:       make(chan IWorker, maxSize),
	}
}

func (t *Troops) SetLogger(logger ILogger) {
	t.Log = logger
}

func (t *Troops) DoJob(f func(args ...interface{}), args ...interface{}) {
	t.JobBuffer <- &Job{exec: f, args: args}
}

func (t *Troops) Run() {
	for i := 0; i < t.MinIdle; i++ {
		s := newSoldier(time.Now().Unix())
		s.setPool(t)
		s.start()
		t.WorkingSize++
	}
	go t.dispach()
}

func (t *Troops) registerWorker(s IWorker) {
	t.Soldiers <- s
}

func (t *Troops) dispach() {
	tick := time.NewTimer(t.MaxIdleTimeout)
	defer tick.Stop()
	for {
		select {
		case job, ok := <-t.JobBuffer:
			if !ok {
				return
			}
			select {
			case worker := <-t.Soldiers:
				worker.Do(job)
			default:
				if t.WorkingSize < t.MaxSize {
					s := newSoldier(time.Now().Unix())
					s.setPool(t)
					s.start()
					t.WorkingSize++
				}
				worker := <-t.Soldiers
				worker.Do(job)
			}
		case <-tick.C:
			if len(t.JobBuffer) >= t.MaxSize {
				break
			}
			reduce := len(t.Soldiers) - t.MinIdle
			for reduce > 0 {
				worker := <-t.Soldiers
				worker.Stop()
				reduce--
				t.WorkingSize--
			}
			tick.Reset(t.MaxIdleTimeout)
		}
	}
}

func (t *Troops) Close() {
	close(t.JobBuffer)
}
