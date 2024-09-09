package comet

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 协议升级
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// func (c *Comet) ServerWebsocket(ch *Channel) {
// 	for {
// 		message, err := ch.conn.ReadMsg()
// 		if ch.done {
// 			return
// 		}
// 		if err != nil {
// 			log.Println("faild to get ws message")
// 			ch.Close()
// 			return
// 		}
// 		log.Println("receive a path message")
// 		f := c.funcRout.Find(message.Path)
// 		if f == nil {
// 			log.Println("wrong method")
// 			continue
// 		}
// 		f(message, ch)
// 		if ch.done {
// 			return
// 		}
// 	}
// }

func (c *Comet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 建立websocket链接
	// log.Println("receive ws connect")
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upGrader fail", err)
		return
	}

	ch := c.NewChannel(ws)

	// add to channels
	c.chmu.Lock()
	c.channels[ch.id] = ch
	c.chmu.Unlock()
	c.serveIO(ch)
}
