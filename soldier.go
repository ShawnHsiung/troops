package troops

import (
	"log"
	"runtime/debug"
)

type IWorker interface {
	Do(*Job)
	Stop()
}

type Soldier struct {
	ID      int64
	JobChan chan *Job
	Pool    *Troops
	Quit    chan struct{}
}

func newSoldier(id int64) *Soldier {
	return &Soldier{
		ID:      id,
		JobChan: make(chan *Job),
		Quit:    make(chan struct{}),
	}
}

func (s *Soldier) setPool(t *Troops) {
	s.Pool = t
}

func (s *Soldier) Do(j *Job) {
	s.JobChan <- j
}

func (s *Soldier) Stop() {
	s.Quit <- struct{}{}
}

func (s *Soldier) start() {
	go func() {

		defer func() {
			if err := recover(); err != nil {
				if s.Pool.Log != nil {
					s.Pool.Log.Errorf("error: %v\nstack: %s", err, string(debug.Stack()))
				} else {
					log.Printf("error: %v\nstack: %s", err, string(debug.Stack()))
				}
			}
		}()

		for {
			s.Pool.registerWorker(s)
			select {
			case job := <-s.JobChan:
				job.exec(job.args...)
			case <-s.Quit:
				return
			}
		}
	}()
}
