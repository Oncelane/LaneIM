package pkg

import (
	"laneIM/src/config"
	"log"

	"github.com/go-redis/redis"
)

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
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("could not connect to redis cluster: %v", err)
	}
	return c
}
