package common

import (
	"hash/fnv"
	"laneIM/src/pkg/laneLog"
	"strconv"
)

type Int64 int64

func (i *Int64) String() string {
	return strconv.FormatInt(int64(*i), 36)
}

func (i *Int64) PasrseString(s string) error {
	t, err := strconv.ParseInt(s, 36, 64)
	*i = Int64(t)
	return err

}

func RedisStrsToInt64(strs []string) ([]int64, error) {
	ret := make([]int64, len(strs))
	var tmp Int64
	for i, str := range strs {
		tmp.PasrseString(str)
		ret[i] = int64(tmp)
	}
	return ret, nil
}

func Int64ToString(in int64) string {
	return strconv.FormatInt(in, 36)
}

func StringTo64(in string) int64 {

	out, err := strconv.ParseInt(in, 36, 64)
	if err != nil {
		laneLog.Logger.Infof("wrong string%s to int64:%v\n", in, err)
		return 404
	}
	return out
}

func HashStringTo64(str string) int64 {
	h := fnv.New64a()
	h.Write([]byte(str))
	return int64(h.Sum64() >> 1)
}
