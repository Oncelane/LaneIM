package pkg_test

import (
	"context"
	"fmt"
	"laneIM/src/config"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"strconv"
	"sync"
	"testing"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func TestSetProduceRedis(t *testing.T) {

	e := pkg.NewEtcd(config.Etcd{Addr: []string{
		"172.29.178.158:51240",
		"172.29.178.158:51241",
		"192.168.141.235:51242"}})
	laneLog.Logger.Infoln("get: ", e.GetAddr("logic"))
	e.SetAddr_WithoutWatch("redis:1", "127.0.0.1:7001")
	e.SetAddr_WithoutWatch("redis:2", "127.0.0.1:7002")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7003")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7004")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7005")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7006")
}

func TestSetLocalRedis(t *testing.T) {

	e := pkg.NewEtcd(config.Etcd{Addr: []string{
		"127.0.0.1:51240",
		"127.0.0.1:51241",
		"127.0.0.1:51242"}})
	laneLog.Logger.Infoln("get: ", e.GetAddr("logic"))
	e.SetAddr_WithoutWatch("redis:1", "127.0.0.1:7001")
	e.SetAddr_WithoutWatch("redis:2", "127.0.0.1:7002")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7003")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7004")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7005")
	e.SetAddr_WithoutWatch("redis:3", "127.0.0.1:7006")
}

func BenchmarkSmembers(b *testing.B) {
	var rd = pkg.NewRedisClient(config.DefaultRedis())
	for range b.N {
		var pipe = rd.Client.Pipeline()
		for i := range 1000 {
			strUserid := strconv.FormatInt(int64(i), 36)
			roomSetKey := fmt.Sprintf("user:%s:roomS", strUserid)
			pipe.SMembers(ctx, roomSetKey)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			laneLog.Logger.Fatalln(err)
		}
		// laneLog.Logger.Infoln("query rows:", len(results))
	}
}

func BenchmarkNoBatchSmembers(b *testing.B) {
	// conn := rd.Client.SlowLog()
	var rd = pkg.NewRedisClient(config.DefaultRedis())
	for range b.N {
		for i := range 1000 {
			strUserid := strconv.FormatInt(int64(i), 36)
			roomSetKey := fmt.Sprintf("user:%s:roomS", strUserid)
			_, err := rd.Client.SMembers(ctx, roomSetKey).Result()
			if err != nil {
				laneLog.Logger.Fatalln(err)
			}
		}
	}
}

func BenchmarkGet(b *testing.B) {
	var rdNoCluster = pkg.NewRedisClientNotCluster(config.DefaultRedis())
	for range b.N {
		for range 1000 {
			// roomSetKey := fmt.Sprintf("user:%s:roomS", strconv.FormatInt(int64(i), 36))

			_, err := rdNoCluster.Client.Get(ctx, "test").Result()
			if err != nil && err != redis.Nil {
				laneLog.Logger.Fatalln(err)
			}
		}
	}
}

func TestRediss(t *testing.T) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:7001",
			"localhost:7002",
			"localhost:7003",
			"localhost:7004",
			"localhost:7005",
			"localhost:7006",
		},
		RouteRandomly: true, // 启用随机路由，以实现负载均衡
	})
	// 记录开始时间
	start := time.Now()
	for i := 0; i < 1000; i++ {
		// 使用客户端发送Redis命令

		err := client.Set(ctx, "name"+strconv.Itoa(i), "张三"+strconv.Itoa(i), 0).Err()
		if err != nil {
			panic(err)
		}

		_, err = client.Get(ctx, "name"+strconv.Itoa(i)).Result()
		if err != nil {
			panic(err)
		}
		// fmt.Println("key", val)
	}
	// 记录结束时间
	end := time.Now()
	// 打印代码运行时长
	fmt.Printf("TestRedis代码运行时长: %v\n", end.Sub(start))
}
func TestGoroutineChannelRedis(t *testing.T) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:7001",
			"localhost:7002",
			"localhost:7003",
			"localhost:7004",
			"localhost:7005",
			"localhost:7006",
		},
		RouteRandomly: true, // 启用随机路由，以实现负载均衡
	})
	// 记录开始时间
	start := time.Now()
	fmt.Println(start)
	var wg sync.WaitGroup
	// 创建一个缓冲channel，用于并发插入数据
	ch1 := make(chan int, 1000)
	// 将100000个值放入channel中
	for i := 0; i < 1000; i++ {
		ch1 <- i
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// 从channel中取出一个值插入到Redis Cluster中
				key := <-ch1
				// 使用客户端发送Redis命令
				err := client.Set(ctx, "name"+strconv.Itoa(key), "张三"+strconv.Itoa(key), 0).Err()
				if err != nil {
					fmt.Println("Error inserting value:", err)
				}
			}
		}()
	}
	close(ch1)
	wg.Wait()
	// 记录结束时间
	end := time.Now()
	fmt.Println(end)
	// 打印代码运行时长
	fmt.Printf("insertGoroutineChannelRedis代码运行时长: %v\n", end.Sub(start))
}
