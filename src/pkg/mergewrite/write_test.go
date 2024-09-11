package mergewrite_test

import (
	"laneIM/src/config"
	"laneIM/src/pkg/laneLog.go"
	"laneIM/src/pkg/mergewrite"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T) {
	w := mergewrite.NewMergeWriter(config.BatchWriter{
		MaxTime:  100,
		MaxCount: 100,
	})

	for seq := range 1000 {
		mockSql := Generator(seq)
		time.Sleep(time.Millisecond)

		go func() {
			_, _ = w.Do("room_comet", mockSql)

			// laneLog.Logger.Infoln(rt, err)
		}()

	}
	select {}
}
func Generator(seq int) func() (any, error) {
	return func() (any, error) {
		//attmept to do sql
		time.Sleep(time.Second)
		laneLog.Logger.Infoln("do sql:", seq)
		return seq, nil
	}
}

func TestWithArg(t *testing.T) {
	w := mergewrite.NewMergeWriter(config.BatchWriter{
		MaxTime:  100,
		MaxCount: 100,
	})

	for seq := range 1000 {
		tmp := seq
		mockSql := GeneratorWithArg(tmp)
		time.Sleep(time.Millisecond)

		go func() {
			w.DoWithArg("room_comet", tmp, mockSql)
		}()

	}
	select {}
}

func GeneratorWithArg(seq int) func(any) (any, error) {
	return func(in any) (any, error) {
		//attmept to do sql
		time.Sleep(time.Second)
		laneLog.Logger.Infoln("input arg:", in.(int))
		return seq, nil
	}
}

func TestGo(t *testing.T) {
	conf := config.BatchWriter{
		MaxTime:  100,
		MaxCount: 10,
	}
	w := mergewrite.NewMergeWriter(conf)
	strSlice := []string{"hello"}
	for i := range 1000 {
		strSlice[0] = strconv.FormatInt(int64(i), 10)
		var tmp = strSlice[0]
		go w.Do(strSlice[0], func() (any, error) {
			laneLog.Logger.Infoln(tmp)
			return nil, nil
		})
	}
	select {}
}
