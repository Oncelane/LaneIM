# discord 大群组高并发分布式 IM

以 discord 万人群组为目标设计，架构设计参考 bilibili 开源 goim，分为 comet（网关），job（数据推送），logic（业务）三大模块，均为无状态节点，模块均可弹性集群部署；集群通讯使用 grpc+protobuf 和 kafka；

comet 网关采用 websocket 长连接，内存池优化消息结构体；logic 业务模块接收 comet 信息，负责消息存储/用户状态查询/网关查询/群成员变动等等；job 推送模块缓存网关信息，转发群消息至 comet；

redis 缓存用户状态，网关路由等信息；kafka 消息队列，解耦业务与网络 IO；laneEtcd 集群服务注册发现；

二级缓存设计，设置 bigcache 一级本地缓存，redis 二级远程缓存；canal 维持 redis 数据一致性；redis 的 sub/pub 维持本地缓存一致性；

scyllaDB 海量消息存储和高效分页拉取; zap+lumberjack 异步日志和日志轮转;singleflight 优化高并发读；请求合并+batch APi 优化高并发读写：以可控较少延迟为代价减少网络请求次数

用户订阅消息机制，对不关心的群聊只接收消息未读数，解决扩散读问题；

离线消息机制：用户上线拉取历史未读信息数目

## 已实现业务

除去基本的加入/创建群组，收发消息外

- 上线后可查询自上次离线以来的未读消息数
- 历史消息分页拉取
- 消息订阅机制，用户对不感兴趣的群组只订阅未读消息数目，感兴趣的群组接收完整消息
  > 业务上所谓“感兴趣”可以定义为是用户点开了页面的群组，也可以是用户设置了特别关心的群组
  > "不感兴趣"是用户没有主动点开的群组
  > 因此还可以额外设置一个"屏蔽"选项，连消息未读数也不接受
  > 实际上 discord 默认对所有群组都是"屏蔽"

# 项目依赖

- [laneEtcd](<(https://github.com/Oncelane/laneEtcd)>) (本人另一个项目，性能相当的 etcd 实现)
- gRpc
- protobuf
- mysql
- scyllaDB
- kafka
- canal
- redis

# 压测一 万人群持续高负载写入测试

压测环境：两台笔记本的 WSL 虚拟机上分别运行一个 comet 实例组成两个 comet 的集群；再将一个 job 和 logic 部署到第二台笔记本上，均连接手机热点组成局域网

压测 case：每个 comet 长连接 5000 个用户，组成 10000 个在线用户的群组，并持续推送消息

| 项目            | 笔记本 1       | 笔记本 2       | 系统之和       |
| --------------- | -------------- | -------------- | -------------- |
| 类型            | 10000 人群推送 | 10000 人群推送 | 10000 人群推送 |
| 推送内容        | "hello"        | "hello"        | "hello"        |
| 持续推送数据量  |                |                | **1340 条/秒** |
| 推送到达        |                |                | **1340 万/秒** |
| 接收流量 (eth0) | 124 KB/s       | 218 KB/s       | 342 KB/s       |
| 发送流量 (eth0) | 98 KB/s        | 298 KB/s       | 396 KB/s       |
| 客户端接收流量  | 32.8MB/s       | 33.4 MB/s      | **66.2MB/s**   |
| 接收流量 (回环) | 185 MB/s       | 231 MB/s       | 416MB/s        |
| 发送流量 (回环) | 185 MB/s       | 231 MB/s       | 416MB/s        |

## 压测截图

server1 中的客户端监控数据

![测试截图1](./docs/压测server2客户端数据.png)

server2 网卡统计数据

![测试截图2](./docs/压测server2网卡.png)

系统消息转发数

![2](./docs/压测每秒转发数.png)

客户端接受流量

![2](./docs/每秒输出接收量.png)

# 压测二 千人群发送消息可用性测试

## 压测环境：

本地单机搭建集群，job，comet，logic 各一个

## 压测 case：

在千人群中，单人发送 N 条消息，并等待所有成员收到所花费的时间

> 注意，由于是单客户端发送消息，因此不是并发，因为每发送一条消息都要得到 comet 的 ack 响应后才能发送下一条。

| 单人发送消息数 | 单人发送 QPS | 总转发数 | 从发送消息到完全接收的总耗时 | 发送方平均每条消息 ACK 延迟 | 等发送后，其他客户端完全接收所有消息耗时 |
| -------------- | ------------ | -------- | ---------------------------- | --------------------------- | ---------------------------------------- |
| 1              | -            | 1000     | 191ms                        | 0.52ms                      | 190ms                                    |
| 36             | -            | 3.6W     | 186ms                        | 0.28ms                      | 177ms                                    |
| 1424           | 3157         | 142.4W   | 556ms                        | 0.31ms                      | 105ms                                    |
| 4.80W          | 3005         | 4808.7W  | 16s                          | 0.33ms                      | 143ms                                    |
| 19.32W         | 3116         | 19322.8W | 62s                          | 0.32ms                      | 152ms                                    |

从消息发送速度来看，单客户端 0.32ms 的消息发送速度（即 ACK 响应速度）是比较理想的，若换成多客户端并发发送，系统每秒处理的速度上线更高，由压测 1 可以看到系统上限在 1000W+ QPS 级别，而本压测中则显示大概在 311 W QPS 级别

并且由于 batch Api 的存在，客户端的多次请求会被合并，模块间的网络通讯次数**并不会跟随请求数的增长而线性增长，对网络 IO 利用率是比较高的**：可以看到即使是 2 亿级别的消息转发量，客户端接收耗时也不会随之增长。（注意，并不是 2 亿条消息只花了 152ms 就全部接收到，而是花了 62s+152ms 就全接收到了。即发送者发送 20 万条消息中的最后一条消息成功后的 152ms ，其余客户端便接收到了）

# 压测三 历史消息连续拉取性能

## 压测 case 1：

群内已存入 1000 条历史消息，1 个客户端每次拉取 20 条历史消息

![历史消息](./docs/每页20拉取1000条历史消息.png)

## 压测 case 2：

群内已存入 1000 条历史消息，1000 个客户端每次拉取 20 条历史消息

## 结论：

由于请求合并会延迟请求，因此即使是单个客户端的请求，每页耗时在 100ms 左右，100ms 用于等待后续相同请求，而实际请求 scyllaDB 耗时仅在 1~2ms。好处是即使当客户端并发请求的数量上升时，所有用户的平均延迟仍然为 100ms 左右

## 补充：

由于 scyllaDB 的时间存储精度只能达到 ms 级别，**对于 1ms 内的重复消息无法准确区分页的边界，只能采取消极应对方法**，使用 <= 的语义拉取历史消息，而不是<

因此最终拉取的分页数目将超过 50，重复的消息将由接收客户端根据 message_id 字段去重

![分页数目](./docs/拉取分页的数目.png)

> 可见，最后总共拉取了 57 次,57\*20 约等于 1140 条数据。
>
> 即在服务端的 scyllaDB 平均每 1ms 存在三条消息的密集度下，客户端不断拉取每页 20 条历史消息，最后得到的消息重复率为 14%。
>
> 证明了虽然我采用的是消极应对策略，带来的额外开销仍在可允许范围内。
>
> 因为实际使用过程中，同一个群聊的消息并发数显然不可能达到如此之高，以至于密集到每一 ms 内都重叠了三四条，因此我大胆估计实际使用中的消息重复率远低于 1%

# 压测 4 - ScyllaDB vs Mysql 的 "count()" 性能对比

## 压测 case:

根据用户的离线时间戳，查询 DB 中属于此用户的未读消息数,求服务器对此请求的处理效率

dao 层执行的核心代码

```sh
# scyllaDB 的 cql 语句：
select  count (1) from %s.messages
where group_id = %d AND timestamp >= '%s'
```

```sh
# mysql 的 sql 语句：
SELECT COUNT(*)
FROM chat_message_mysqls
WHERE group_id = ? AND timestamp >= ?;
```

补充阅读: mysql 中 count(\*)和其他 count 方法的性能对比

> count（主键）：InnoDB 引擎会遍历整张表，把每一行的 主键 id 值都取出来，返回给服务层。服务层拿到主键后，直接按行进行累加（主键不可能为 null）
> count（字段）：没有 not null 约束，InnoDB 引擎会遍历整张表把每一行的字段值都取出来，返回给服务层，服务层判断是否为 null，不为 null，计数累加。有 not null 约束，InnoDB 引擎会遍历整张表把每一行的字段值都取出来，返回给服务层，直接按行进行累加。
> count（1）：InnoDB 引擎遍历整张表，但不取值。服务层对于返回的每一行，放一个数字”1“进行，直接按行进行累加。
> count（\*）：InnoDB 引擎并不会把全部字段取出来，而是专门做了优化。不取值，服务层直接按行累加。
>
> 按照效率排序的话，count（字段）< count（主键 id）< count（1）约等于 count（\*），所以尽量使用 count（\*）。
> [原文链接](https://blog.csdn.net/m0_65152767/article/details/140052693)

## 压测性能：

| 服务器 QPS | 不使用协程池做并发控制     | 使用协程池，控制最大并发数为 150 |
| ---------- | -------------------------- | -------------------------------- |
| Mysql      | Error：too many connection | 74.616μs/op                      |
| ScyallDB   | 2.756μs/op                 | 92.995μs/op                      |

## 压测截图：

不进行并发限制的 Mysql 直接触发 too many connection 错误，并发连接数超过 141 即爆炸,应该无需截图证明了

不进行并发限制的 scyllaDB 查询性能
![不限制的](./docs/不进行pool的scyllaDB查询性能.png)

> ScyllaDB 似乎自带并发控制，裸并发 566W 个协程也不存在问题（当然，即使 566W 个协程是在同一个循环内快速启动的，实际每一刻的协程数目应该远低于 566W 个，但也至少比 150 高得多）

限制并发数目为 150 的 ScyllaDB：
![限制并发数目的ScyllaDB](./docs/限制150pool的scyllaDB性能.png)

**结论 1：scyllaDB 貌似自带并发限制功能，直接开几万个协程并发访问也不会存在问题，且性能比自我限制要高得多**

限制并发数目为 150 的 Mysql：
![限制并发数目的Mysql](./docs/进行150pool的mysql性能.png)

# 为什么选择 scyllaDB 存储历史消息

1.  业务需求：历史消息没有太多跨表 join 的需求，即使被独立开来使用 nosql 存储，也几乎不会影响到相关业务

2.  性能需求：noSql 在数据的读，写都远快于 sql，上有压测为证

3.  分布式需求：传统 sql 注重事务和强一致性，**属于 CA 模型**，若要部署为集群，则要采用手动分片，也就是分库分表的方式做分布式存储。由此引发一系列的问题比如数据迁移，数据的负载均衡，都需要使用者自己解决。**而 scyllaDB 支持去中心化的集群部署**，数据使用 partition key 做负载均衡。好处是：一方面 partition key 的 value 值掌握在使用者手里，另一方面，麻烦的数据迁移和负载均衡都由 scyllaDB 自己完成，更适合分布式海量存储。

# 编译

```sh
make all -i

# 也可以分别编译
make build-comet
make build-job
make build-logic
```

修改 makefile 中的集群数目

```sh
# Default number of clusters

# local Cluster
Ncomet ?= 2
Njob ?= 1
Nlogic ?= 1

# network Cluster
Ncometp ?= 1
Njobp ?= 1
Nlogicp ?= 1

```

# 启动：

首先要检查有无启动 kafka，canal，redis，laneEtcd（请查看[此项目地址](https://github.com/Oncelane/laneEtcd)）

本地集群方式（推荐）:

打开三个终端，分别运行三个命令以启动三个集群

logic 集群

```sh
make run-logic -i
```

comet 集群

```sh
make run-comet -i
```

job 集群

```sh
make run-job -i
```

canal 服务器是可选的，它可以用于保障 sql 和 redis 之间，以及 redis 和 localcache 之间的数据一致性，但是有一定的延迟

```sh
cd src/cmd/canal/
go run main.go
```

# 测试

```sh
# 进入src/cmd
cd src/client
# 使用vscode之类的ide手动运行测试项目即可
```

# 依赖安装

安装 protobuf grpc complier
sudo apt-get install protobuf-compiler

gRPC Go 插件

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

# protoc 命令

```sh
protoc --go_out=.. --go-grpc_out=.. --go-grpc_opt=require_unimplemented_servers=false -I. -Iproto proto/msg/msg.proto proto/comet/comet.proto proto/logic/logic.proto
```

# scylla 集群

## 安装：

[官网网址](https://www.scylladb.com/download/?platform=ubuntu-20.04&version=scylla-6.1#open-source)

可能需要科学上网

init address

```sh
sudo mkdir -p /etc/apt/keyrings
sudo gpg --homedir /tmp --no-default-keyring --keyring /etc/apt/keyrings/scylladb.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 491c93b9de7496a7
```

```sh
sudo curl -L --output /etc/apt/sources.list.d/scylla.list https://downloads.scylladb.com/deb/ubuntu/scylla-6.1.list
```

Install packages

```sh
sudo apt-get update
sudo apt-get install -y scylla
```

Set Java to 1.8 release

```sh
sudo apt-get update
sudo apt-get install -y openjdk-8-jre-headless
sudo update-java-alternatives --jre-headless -s java-1.8.0-openjdk-amd64
```

## 配置:

setup

```sh
# scylla_setup accepts command line arguments as well! For easily provisioning in a similar environment than this, type:

    scylla_setup --no-raid-setup --online-discard 1 --nic eth0 \
                 --io-setup 1 --no-memory-setup --no-rsyslog-setup

# Also, to avoid the time-consuming I/O tuning you can add --no-io-setup and copy the contents of /etc/scylla.d/io*
# Only do that if you are moving the files into machines with the exact same hardware
```

config

```sh
sudo nano /etc/scylla/scylla.yaml

# cluster_name: 'laneIM'
```

set ScyllaDB to developer mode.

```sh
sudo scylla_dev_mode_setup --developer-mode 1
```

## 启动

```sh
sudo systemctl start scylla-server.service
```

check status

```sh
sudo systemctl status scylla-server.service
```

# redis 集群

## 搭建

启动六个单独节点

```sh
bash redisClusterStart.sh
```

组建集群

```sh
redis-cli --cluster create 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 127.0.0.1:7006 --cluster-replicas 1
```

## 下次启动

```sh
# 启动
bash redisClusterStart.sh
# 关闭
bash redisClusterClose.sh
```

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

# kafka 集群

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

列出所有键：
使用 etcdctl 列出所有键：

```bash
etcdctl get "" --prefix --keys-only
```

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
# 先启动logic会自动建表
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
canal.instance.dbPassword=QTLVb6BaeeaJsFMT
# debian-sys-maint
# QTLVb6BaeeaJsFMT
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

```sh
#设置RoomMgr信息：
HMSET room:{RoomID} OnlineCount 10
SADD room:{RoomID}:users_set {UserID}
SADD room:{RoomID}:comets_set {CometAddr}

# 设置UserMgr信息：
HMSET user:{UserID} CometAddr {CometAddr} Online true
SADD user:{UserID}:rooms_set {RoomID}

#设置CometMgr信息：
SADD comet:{CometAddr}:rooms_set {RoomID}
```

# next
