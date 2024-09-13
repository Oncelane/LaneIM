package batch_test

import (
	"laneIM/src/pkg/batch"
	"laneIM/src/pkg/laneLog"
	"testing"
	"time"
)

func TestBatch(t *testing.T) {
	b := batch.NewBatchArgs[int](15, time.Millisecond*100, func(in []*int) {
		for i := range in {
			laneLog.Logger.Infof("%d ", *in[i])
		}
		laneLog.Logger.Infoln("=====end=====")
	})
	b.Start()
	for i := range 100 {
		time.Sleep(time.Millisecond * 50)
		t := i
		b.Add(&t)
	}
	select {}
}
