package pkg

import (
	"log"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Addrs  []string
	Client *redis.ClusterClient
}

func NewRedisClient(addrs []string) *RedisClient {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs, // Redis 节点地址
		// 如果你的集群需要身份验证，请提供密码
		// Password: "yourpassword",
	})
	c := &RedisClient{
		Addrs:  addrs,
		Client: rdb,
	}
	// 测试连接
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("could not connect to redis cluster: %v", err)
	}
	return c
}
