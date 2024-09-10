package localCache

import (
	"strconv"
	"strings"

	"github.com/allegro/bigcache"
)

func UserComet(cache *bigcache.BigCache, userid int64) (string, error) {
	key := "user:comet" + strconv.FormatInt(userid, 36)
	data, err := cache.Get(key)
	if err != nil {
		// log.Println("miss local cache user:comet", err)
		return "", err
	}
	return string(data), nil
}

func SetUserComet(cache *bigcache.BigCache, userid int64, comet string) error {
	key := "user:comet" + strconv.FormatInt(userid, 36)
	return cache.Set(key, []byte(comet))
}

func DelUserComet(cache *bigcache.BigCache, userid int64) error {
	key := "user:comet" + strconv.FormatInt(userid, 36)
	return cache.Delete(key)
}

func UserRoomid(cache *bigcache.BigCache, userid int64) ([]int64, error) {
	key := "user:room" + strconv.FormatInt(userid, 36)
	data, err := cache.Get(key)
	if err != nil {
		// log.Println("miss local cache user:user", err)
		return nil, err
	}
	rawInt64 := strings.Split(string(data), ";")
	userids := make([]int64, len(rawInt64))
	for i := range rawInt64 {
		userid, err := strconv.ParseInt(rawInt64[i], 36, 64)
		if err != nil {
			return nil, err
		}
		userids[i] = userid
	}
	return userids, nil
}

func SetUserRoomid(cache *bigcache.BigCache, userid int64, roomids []int64) error {
	key := "user:room" + strconv.FormatInt(userid, 36)
	rawroomids := make([]string, len(roomids))
	for i := range roomids {
		raw := strconv.FormatInt(roomids[i], 36)
		rawroomids[i] = raw
	}
	return cache.Set(key, []byte(strings.Join(rawroomids, ";")))
}

func DelUserRoomid(cache *bigcache.BigCache, userid int64) error {
	key := "user:room" + strconv.FormatInt(userid, 36)
	return cache.Delete(key)
}
