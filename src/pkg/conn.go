package pkg

import (
	"laneIM/proto/msg"
	"laneIM/src/pkg/laneLog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type MsgReadWriteCloser interface {
	ReadMsg() (message *msg.MsgBatch, err error)
	WriteMsg(message *msg.MsgBatch) error
	PassiceClose() error
	ForceClose() error
	PoolPut(message *msg.MsgBatch)
}

type ConnWs struct {
	conn *websocket.Conn
	pool *MsgPool
	once sync.Once
}

var _ MsgReadWriteCloser = new(ConnWs)

func NewConnWs(conn *websocket.Conn, pool *MsgPool) *ConnWs {
	return &ConnWs{
		conn: conn,
		pool: pool,
	}
}

func (w *ConnWs) ReadMsg() (message *msg.MsgBatch, err error) {
	mtype, p, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if mtype == websocket.CloseMessage {
		w.conn.Close()
		return nil, &websocket.CloseError{Code: 1000}
	}
	// message = w.pool.Get()
	message = new(msg.MsgBatch)
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

func (w *ConnWs) PassiceClose() error {
	return w.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*3))
}

func (w *ConnWs) ForceClose() error {
	var err error = nil
	w.once.Do(func() { err = w.conn.Close() })
	return err
}

func (w *ConnWs) PoolPut(message *msg.MsgBatch) {
	w.pool.Put(message)
}
