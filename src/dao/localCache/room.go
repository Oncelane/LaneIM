package localCache

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/allegro/bigcache"
)

func Cache(time time.Duration) *bigcache.BigCache {
	rt, err := bigcache.NewBigCache(bigcache.DefaultConfig(time))
	if err != nil {
		log.Panicln("faild to create local cache")
	}
	return rt
}

func RoomComet(cache *bigcache.BigCache, roomid int64) ([]string, error) {
	data, err := cache.Get(strconv.FormatInt(roomid, 36))
	if err != nil {
		// log.Println("miss local cache room:comet", err)
		return nil, err
	}
	return strings.Split(string(data), ";"), nil
}

func SetRoomComet(cache *bigcache.BigCache, roomid int64, comets []string) error {
	return cache.Set(strconv.FormatInt(roomid, 36), []byte(strings.Join(comets, ";")))
}

func DelRoomComet(cache *bigcache.BigCache, roomid int64) error {
	return cache.Delete("room:comet" + strconv.FormatInt(roomid, 36))
}

func RoomUserid(cache *bigcache.BigCache, roomid int64) ([]int64, error) {
	data, err := cache.Get(strconv.FormatInt(roomid, 36))
	if err != nil {
		// log.Println("miss local cache room:user", err)
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

func SetRoomUserid(cache *bigcache.BigCache, roomid int64, userids []int64) error {
	rawuserids := make([]string, len(userids))
	for i := range userids {
		raw := strconv.FormatInt(userids[i], 36)
		rawuserids[i] = raw
	}
	return cache.Set(strconv.FormatInt(roomid, 36), []byte(strings.Join(rawuserids, ";")))
}

func DelRoomUserid(cache *bigcache.BigCache, roomid int64) error {
	return cache.Delete("room:user" + strconv.FormatInt(roomid, 36))
}
