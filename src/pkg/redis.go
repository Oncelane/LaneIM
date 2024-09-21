package pkg

import (
	"context"
	"laneIM/src/config"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	Addrs  []string
	Client *redis.ClusterClient
}

func NewRedisClient(conf config.Redis) *RedisClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: conf.Addr, // Redis 节点地址
		// 如果你的集群需要身份验证，请提供密码
		// Password: "yourpassword",
	})
	c := &RedisClient{
		Addrs:  conf.Addr,
		Client: rdb,
	}
	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to redis cluster: %v", err)
	}
	return c
}

type RedisClientNotCluster struct {
	Addrs  []string
	Client *redis.Client
}

func NewRedisClientNotCluster(conf config.Redis) *RedisClientNotCluster {
	rdb := redis.NewClient(&redis.Options{
		Addr: conf.Addr[1], // Redis 节点地址
		// 如果你的集群需要身份验证，请提供密码
		// Password: "yourpassword",
	})
	c := &RedisClientNotCluster{
		Addrs:  conf.Addr,
		Client: rdb,
	}
	// 测试连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to redis cluster: %v", err)
	}
	return c
}
