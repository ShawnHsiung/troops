package troops

import (
	"testing"
	"time"
)

func TestDo(t *testing.T) {

	tps := NewTroops(100, 10)

	f := func(i ...interface{}) {
		t.Logf("%v exit", i)
	}

	tps.Run()

	for i := 0; i < 10; i++ {
		tps.DoJob(f, i)
	}

	time.Sleep(2 * time.Second)
}
