package localCache

import (
	"strconv"
	"strings"

	"laneIM/src/pkg/laneLog.go"

	"github.com/allegro/bigcache"
)

func UserComet(cache *bigcache.BigCache, userid int64) (string, error) {
	key := "user:comet" + strconv.FormatInt(userid, 36)
	data, err := cache.Get(key)
	if err != nil {
		// laneLog.Logger.Infoln("miss local cache user:comet", err)
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
		// laneLog.Logger.Infoln("miss local cache user:user", err)
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

func UserRoomidBatch(cache *bigcache.BigCache, userids []int64) ([][]int64, bool) {
	rt := make([][]int64, len(userids))
	full := true
	for i := range userids {
		key := "user:room" + strconv.FormatInt(userids[i], 36)
		data, err := cache.Get(key)
		if err != nil {
			// laneLog.Logger.Infoln("miss local cache user:user", err)
			full = false
			continue
		}
		// laneLog.Logger.Debugf("localcache get userroom key[%s] value[%s]", key, string(data))
		rawInt64 := strings.Split(string(data), ";")
		roomids := make([]int64, len(rawInt64))
		if len(rawInt64) == 0 {
			full = false
			continue
		}
		for i := range rawInt64 {
			roomid, err := strconv.ParseInt(rawInt64[i], 36, 64)
			if err != nil {
				laneLog.Logger.Errorln("faild parse roomid", err)
				full = false
				continue
			}
			roomids[i] = roomid
		}
		rt[i] = roomids
	}
	return rt, full
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

func SetUserRoomidBatch(cache *bigcache.BigCache, userids []int64, roomidss [][]int64) error {
	for i := range userids {
		key := "user:room" + strconv.FormatInt(userids[i], 36)
		rawroomids := make([]string, len(roomidss[i]))
		for j := range roomidss[i] {
			raw := strconv.FormatInt(roomidss[i][j], 36)
			rawroomids[j] = raw
		}
		// laneLog.Logger.Debugf("localcache set key[%s],value[%v]", key, strings.Join(rawroomids, ";"))
		err := cache.Set(key, []byte(strings.Join(rawroomids, ";")))
		if err != nil {
			laneLog.Logger.Errorln("batch set user roomid localcache faild", err)
		}
	}
	return nil
}

func DelUserRoomid(cache *bigcache.BigCache, userid int64) error {
	key := "user:room" + strconv.FormatInt(userid, 36)
	return cache.Delete(key)
}
