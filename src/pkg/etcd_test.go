package pkg_test

import (
	"laneIM/src/pkg"
	"log"
	"testing"
)

func TestEtcd(t *testing.T) {
	e := pkg.NewEtcd()
	log.Println("userid 1: ", e.GetUserRoomid(1))

	log.Println(e.AddUserRoomid(1, 2345))

	log.Println("userid 1: ", e.GetUserRoomid(1))
	log.Println(e.SetUserRoomid(1, []int64{1, 2, 3, 4, 5, 6}))
	log.Println("userid 1: ", e.GetUserRoomid(1))

}
