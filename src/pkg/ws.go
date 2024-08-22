package pkg

import (
	"laneIM/proto/msg"
	"log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Ws struct {
	conn *websocket.Conn
	pool *MsgPool
}

type RWMsg interface {
	ReadMsg() (message *msg.Msg, err error)
	WriteMsg(message *msg.Msg) error
}

func NewWs(conn *websocket.Conn, pool *MsgPool) *Ws {
	return &Ws{
		conn: conn,
		pool: pool,
	}
}

func (w *Ws) ReadMsg() (message *msg.Msg, err error) {
	_, p, err := w.conn.ReadMessage()
	if err != nil {
		log.Println("read err")
		return nil, err
	}
	message = w.pool.Get()
	err = proto.Unmarshal(p, message)
	if err != nil {
		log.Println("unmarshal err", err)
		return nil, err
	}
	return
}

func (w *Ws) WriteMsg(message *msg.Msg) error {
	p, err := proto.Marshal(message)
	if err != nil {
		log.Println("marshal err", err)
		return err
	}
	err = w.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		log.Println("websocket write err", err)
		return err
	}
	return nil
}
