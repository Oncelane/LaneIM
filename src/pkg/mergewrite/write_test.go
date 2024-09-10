package mergewrite_test

import (
	"laneIM/src/config"
	"laneIM/src/pkg/mergewrite"
	"log"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T) {
	w := mergewrite.NewMergeWriter(config.BatchWriter{
		MaxTime:  100,
		MaxCount: 100,
	})

	for seq := range 1 {
		mockSql := Generator(seq)
		time.Sleep(time.Millisecond)

		go func() {
			rt, err := w.Do("room_comet", mockSql)

			log.Println(rt, err)
		}()

	}
	select {}
}

func Generator(seq int) func() (any, error) {
	return func() (any, error) {
		//attmept to do sql
		time.Sleep(time.Second)
		log.Println("do sql:", seq)
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
			log.Println(tmp)
			return nil, nil
		})
	}
	select {}
}
