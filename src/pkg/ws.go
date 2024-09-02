package pkg

import (
	"laneIM/proto/msg"
	"log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type MsgReadWriteCloser interface {
	ReadMsg() (message *msg.Msg, err error)
	WriteMsg(message *msg.Msg) error
	Close() error
}
type ConnWs struct {
	conn *websocket.Conn
	pool *MsgPool
}

var _ MsgReadWriteCloser = new(ConnWs)

func NewConnWs(conn *websocket.Conn, pool *MsgPool) *ConnWs {
	return &ConnWs{
		conn: conn,
		pool: pool,
	}
}

func (w *ConnWs) ReadMsg() (message *msg.Msg, err error) {
	_, p, err := w.conn.ReadMessage()
	if err != nil {
		log.Println("read err:", err)
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

func (w *ConnWs) WriteMsg(message *msg.Msg) error {
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

func (w *ConnWs) Close() error {
	return w.conn.Close()
}
