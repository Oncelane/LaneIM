package pkg

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg/laneLog"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type MsgReadWriteCloser interface {
	ReadMsg() (message *msg.MsgBatch, err error)
	WriteMsg(message *msg.MsgBatch) error
	Close() error
	ForceClose() error
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

func (w *ConnWs) ReadMsg() (message *msg.MsgBatch, err error) {
	_, p, err := w.conn.ReadMessage()
	if err != nil {
		// laneLog.Logger.Infoln("read err:", err)
		return nil, err
	}
	message = w.pool.Get()
	err = proto.Unmarshal(p, message)
	if err != nil {
		laneLog.Logger.Infoln("unmarshal err", err)
		return nil, err
	}
	return
}

func (w *ConnWs) WriteMsg(message *msg.MsgBatch) error {
	p, err := proto.Marshal(message)
	if err != nil {
		laneLog.Logger.Infoln("marshal err", err)
		return err
	}
	err = w.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		laneLog.Logger.Infoln("websocket write err", err)
		return err
	}
	return nil
}

func (w *ConnWs) Close() error {
	return w.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*3))
}

func (w *ConnWs) ForceClose() error {
	return w.conn.Close()
}
