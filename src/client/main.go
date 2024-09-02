package main

import (
	"encoding/json"
	"laneIM/proto/msg"
	"laneIM/src/comet"
	"laneIM/src/pkg"
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	service := "ws://localhost:40051/ws" // 替换为你的WebSocket服务地址

	// 连接到WebSocket服务
	conn, _, err := websocket.DefaultDialer.Dial(service, nil)
	if err != nil {
		log.Fatal("连接错误:", err)
	}
	defer conn.Close()
	log.Println("连接到comet:", service)
	connMsg := pkg.NewConnWs(conn, pkg.NewMsgPool())
	// 发送消息到服务器
	token := comet.Token{
		Content: []byte("this is 21 token"),
		Userid:  21,
	}
	tokenDate, err := json.Marshal(token)
	if err != nil {
		log.Println("faild to encode token")
	}
	log.Println("授权，并假定成功过")
	err = connMsg.WriteMsg(&msg.Msg{
		Path: "Auth",
		Seq:  1,
		Byte: tokenDate,
	})

	if err != nil {
		log.Println("发送消息错误:", err)
	}
	log.Println("发送房间消息")
	err = connMsg.WriteMsg(&msg.Msg{
		Path: "room",
		Seq:  1,
		Byte: []byte("hello form 21 room message"),
	})
	if err != nil {
		log.Println("发送消息错误:", err)
	}
	// 读取服务器响应
	go func() {
		for {
			message, err := connMsg.ReadMsg()
			if err != nil {
				log.Println("接收消息错误:", err)
				return
			}
			log.Printf("接收到消息: %s", message.String())
		}
	}()
	select {}
}
