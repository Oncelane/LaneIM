package comet

import (
	"laneIM/proto/msg"
	"log"
)

type WsHandler func(in *msg.Msg, ch *Channel)

type WsFuncRouter struct {
	fmap map[string]WsHandler
}

func NewWsFuncRouter() *WsFuncRouter {
	return &WsFuncRouter{
		fmap: make(map[string]WsHandler),
	}
}

func (w *WsFuncRouter) Use(path string, f WsHandler) {
	w.fmap[path] = f
	log.Println("registe method:", path)
}

func (w *WsFuncRouter) Find(path string) WsHandler {
	if rt, exist := w.fmap[path]; exist {
		return rt
	}
	log.Println("faild to find method:", path)
	return nil
}
