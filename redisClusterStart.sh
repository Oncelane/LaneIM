#!/bin/bash

# 启动 Redis 节点
for port in {7001..7006}
do
  echo "Starting Redis node on port $port"
  cd redisConf/$port && redis-server ./redis.conf &
  cd - > /dev/null
done
