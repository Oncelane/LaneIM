# laneIM golang 分布式 im

安装 protobuf grpc complier
sudo apt-get install protobuf-compiler

gRPC Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc 命令
protoc --go_out=.. --go-grpc_out=.. --go-grpc_opt=require_unimplemented_servers=false -I. -Iproto proto/msg/msg.proto proto/comet/comet.proto proto/logic/logic.proto

# 启动 redis 集群

redis-server ./redis.conf

redis-cli --cluster create 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 127.0.0.1:7006 --cluster-replicas 1

​ 执行 redis-cli -p 7001 进入客户端并通过 `info replication` 查看集群信息

通过`cluster nodes`查看集群关系

通过`cluster info`查看集群信息

# kafka 客户端

客户端：confluent-kafka-go

下载 lib

```bash
git clone https://github.com/edenhill/librdkafka.git
cd librdkafka
```

安装

```bash
./configure --prefix /usr
make
sudo make install
```

配置

```bash
export PKG_CONFIG_PATH=/usr/lib/pkgconfig
```

# 启动 kafka 集群

确保 Kafka 版本支持 Kraft 模式（Kafka 2.8.0 及以上版本）
下载

```sh
   wget https://downloads.apache.org/kafka/3.7.1/kafka_2.12-3.7.1.tgz
   tar -xzf kafka_2.13-3.0.0.tgz
   cd kafka_2.13-3.0.0
```

编辑 server.properties 配置文件
在 config 目录下，打开 server.properties 文件，并添加或修改以下配置：

server.properties

```sh
   # Set the process roles to broker and controller
   process.roles=broker,controller

   # Specify the controller listener name
   controller.listener.names=CONTROLLER

   # Define listeners for broker and controller
   listeners=PLAINTEXT://localhost:9092,CONTROLLER://localhost:9093

   # Define inter-broker listener name
   inter.broker.listener.name=PLAINTEXT

   # Specify the log directory for metadata
   log.dirs=/tmp/kraft-combined-logs

   # Specify the metadata quorum
   controller.quorum.voters=0@localhost:9093

   # Cluster ID
   broker.id=0

   # Initial cluster ID (generate a new cluster ID if this is the first time you're starting this cluster)
   # You can generate a new cluster ID using the kafka-storage.sh tool
```

生成集群 ID
使用 Kafka 提供的工具 kafka-storage.sh 生成一个新的集群 ID，并格式化日志目录。

```sh
 bin/kafka-storage.sh random-uuid
```

将生成的 UUID 替换到以下命令中：

```sh
   bin/kafka-storage.sh format -t <生成的UUID> -c config/server.properties
```

二、启动 Kafka
启动 Kafka 服务器

```sh
   bin/kafka-server-start.sh config/server.properties
```

三、验证 Kafka 是否在 Kraft 模式下运行
检查日志
查看 Kafka 日志，确保没有错误，并且 Kafka 正确地启动了控制器和代理。
使用 Kafka 客户端
使用 Kafka 客户端来创建主题、生产和消费消息，确保 Kafka 正常运行。

```sh

   # 创建一个主题
   bin/kafka-topics.sh --create --topic test-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

   # 生产消息
   bin/kafka-console-producer.sh --topic test-topic --bootstrap-server localhost:9092

   # 消费消息
   bin/kafka-console-consumer.sh --topic test-topic --from-beginning --bootstrap-server localhost:9092
```

# Room design
in redis

roomMgr set int

room:online:%id int

room:comet:%id set string

room:userid:%id set int

user:comet:%id string

user:room:%id set string

user:online:%id bool

room{
   new
   del

   joinUserid
   quitUserid
   queryUserid
   onlienUser
   offlineUser
}

user{
   new
   del

   online
   offline
   queryOnline
}





