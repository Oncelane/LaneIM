package common

import "strconv"

type Int64 int64

func (i *Int64) String() string {
	return strconv.FormatInt(int64(*i), 36)
}

func (i *Int64) PasrseString(s string) error {
	t, err := strconv.ParseInt(s, 36, 64)
	*i = Int64(t)
	return err

}
