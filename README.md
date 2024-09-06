# laneIM

集群部署的分布式 IM，主要使用到的组件：etcd，kafka，redis，canal，mysql，grpc，protobuf

分为三个模块，comet 集群（网关/代理）， job 集群（消息推送），logic（业务服务器：登录，上下线）
mysql+canel+redis 集群（用户状态，房间信息，路由）
kafka 集群 消息推送队列
etcd 集群 服务注册发现（后续换成手搓的 raft）
grpc+protobuf 微服务通讯

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

查看所有 key

```bash

```

删除所有 kv

```bash

```

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

依赖 jdk11

```bash
export JAVA_HOME=/usr/local/java/jdk1.8.0_411
export JRE_HOME=${JAVA_HOME}/jre
export CLASSPATH=.:${JAVA_HOME}/lib:${JRE_HOME}/lib
export PATH=${JAVA_HOME}/bin:$PATH
```

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

# etcd 集群

goreman 工具

```sh
go install github.com/mattn/goreman@latest
goreman -f local-cluster-profile start
```

下载 etcd

```sh
wget https://github.com/etcd-io/etcd/releases/download/v3.5.15/etcd-v3.5.15-linux-amd64.tar.gz
cd etcd
./build.sh
nano ~/.bashrc
export PATH="$PATH:$GOPATH/src/github.com/etcd-io/etcd/bin"
source ~/.bashrc
```

删除所有键

列出所有键：
使用 etcdctl 列出所有键：

```bash
etcdctl get "" --prefix --keys-only
```

这会列出所有键。

删除所有键：
使用 etcdctl del 命令删除所有键：

```bash
etcdctl del "" --prefix
--prefix 选项会删除以指定前缀开头的所有键。由于指定了空字符串 "" 作为前缀，这会删除所有键。
```

# mysql 安装

```bash
sudo apt update
sudo apt install mysql-server
# 查看密码
sudo cat /etc/mysql/debian.cnf
```

# mysql 建库

```bash
# 创建数据库laneIM
# 运行room_test.go
```

# canal 安装

```bash
sudo mkdir -p /opt/canal
sudo chmod 777 /opt/canal && cd /opt/canal
wget https://github.com/alibaba/canal/releases/download/canal-1.1.7/canal.deployer-1.1.7.tar.gz
tar -zxvf canal.deployer-1.1.7.tar.gz
```

# canal 配置

```bash
mkdir -p conf/laneIM
cp conf/example/instance.properties conf/laneIM/

```

```sh
nano conf/canal.properties
#################################################
#########               destinations            #############
#################################################
canal.destinations = laneIM
```

```sh
nano conf/laneIM/instance.properties
#################################################
## mysql serverId , v1.0.26+ will autoGen
# canal.instance.mysql.slaveId=0

# enable gtid use true/false
canal.instance.gtidon=false

# position info
canal.instance.master.address=127.0.0.1:3306
canal.instance.master.journal.name=
canal.instance.master.position=
canal.instance.master.timestamp=
canal.instance.master.gtid=

# rds oss binlog
canal.instance.rds.accesskey=
canal.instance.rds.secretkey=
canal.instance.rds.instanceId=

# table meta tsdb info
canal.instance.tsdb.enable=true
#canal.instance.tsdb.url=jdbc:mysql://127.0.0.1:3306/canal_tsdb
#canal.instance.tsdb.dbUsername=canal
#canal.instance.tsdb.dbPassword=canal

#canal.instance.standby.address =
#canal.instance.standby.journal.name =
#canal.instance.standby.position =
#canal.instance.standby.timestamp =
#canal.instance.standby.gtid=

# username/password
canal.instance.dbUsername=debian-sys-maint
canal.instance.dbPassword=FJho5xokpFqZygL5
canal.instance.connectionCharset = UTF-8
# enable druid Decrypt database password
canal.instance.enableDruid=false
#canal.instance.pwdPublicKey=MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALK4BUxdDltRRE5/zXpVEVPUgunvscYFtEip3pmLlhrWpacX7y7GCMo2/JM6LeHmiiNdH1FWgGCpUfircSwlWKUCAwEAAQ==

# table regex
canal.instance.filter.regex=.*\\..*
# table black regex
canal.instance.filter.black.regex=mysql\\.slave_.*
# table field filter(format: schema1.tableName1:field1/field2,schema2.tableName2:field1/field2)
#canal.instance.filter.field=test1.t_product:id/subject/keywords,test2.t_company:id/name/contact/ch
# table field black filter(format: schema1.tableName1:field1/field2,schema2.tableName2:field1/field2)
#canal.instance.filter.black.field=test1.t_product:subject/product_image,test2.t_company:id/name/contact/ch

# mq config
canal.mq.topic=laneIM
# dynamic topic route by schema or table regex
#canal.mq.dynamicTopic=mytest1.user,topic2:mytest2\\..*,.*\\..*
canal.mq.partition=0
# hash partition config
#canal.mq.enableDynamicQueuePartition=false
#canal.mq.partitionsNum=3
#canal.mq.dynamicTopicPartitionNum=test.*:4,mycanal:6
#canal.mq.partitionHash=test.table:id^name,.*\\..*
#
# multi stream for polardbx
canal.instance.multi.stream.on=false
#################################################
```

# canal 启动

```sh
bash bin/startup.sh
```

# Room design

in redis

room:mgr set int64

user:mgr set int64

room:online:%id int

room:comet:%id set string

room:userid:%id set int64

user:comet:%id string

user:room:%id set string

user:online:%id bool

```
Allroomid

room{
   new
   del

   joinUserid
   quitUserid
   queryUserid

   putComet
   delComet
   queryComet

   countUserid
}

Alluserid

user{
   new
   del

   online
   offline
   queryOnline

   joinRoomid
   quitRoomid
}
```

# next

## comet

room 模块
websocekt 长连接模块

job 集群部署
comet 集群部署
