#!/bin/bash

# 停止 Redis 节点
for port in {7001..7006}
do
  echo "Stopping Redis node on port $port"
  pid=$(ps -ef | grep "redis-server .*:$port" | grep -v grep | awk '{print $2}')
  if [ -n "$pid" ]; then
    kill $pid
  fi
done
