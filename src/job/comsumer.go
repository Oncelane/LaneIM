package job

import (
	"laneIM/proto/comet"
	"laneIM/proto/logic"
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type MyConsumer struct {
	// 可以添加其他字段，如日志记录器、统计信息等
	job *Job
}

func (consumer *MyConsumer) Setup(session sarama.ConsumerGroupSession) error {
	// 在这里可以执行一些初始化操作，比如日志记录
	return nil
}

func (consumer *MyConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	// 在这里可以执行一些清理操作
	return nil
}

func (consumer *MyConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// log.Printf("Received message: from %s/%d\n", message.Topic, message.Partition)
		// 处理消息...
		protoMsg := &logic.SendMsgReq{}
		err := proto.Unmarshal(message.Value, protoMsg)
		if err != nil {
			log.Println("wrong protobuf decode")
		}
		switch protoMsg.Path {
		case "sendRoom":
			// 可能的初始化room，获取room所处的comets
			room := consumer.job.Bucket(protoMsg.Roomid).GetRoom(protoMsg.Roomid)
			// for addr := range room.info.Server {
			log.Printf("message to roomid:%d hava comet:%v", room.roomid, room.info.Server)
			// 检查是否实际连接上了comet
			consumer.job.Push(&comet.RoomReq{
				Roomid: protoMsg.Roomid,
				Data:   protoMsg.Data,
			})

			// }
		}
		// 标记消息为已消费
		session.MarkMessage(message, "")
	}
	return nil
}
