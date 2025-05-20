package utils

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/user2083251241/ebidsystem/internal/app/config"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
	redisOnce   sync.Once
	redisMutex  sync.Mutex //互斥锁
)

func InitRedis() {
	redisOnce.Do(func() {
		redisMutex.Lock()
		defer redisMutex.Unlock()
		addr := config.Get("REDIS_ADDR")
		password := config.Get("REDIS_PASSWORD")
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})
		//测试连接是否成功：
		_, err := RedisClient.Ping(Ctx).Result()
		if err != nil {
			log.Fatalf("Redis 连接失败: %v", err) //终止应用并记录错误
		}
		log.Println("Redis 连接成功")
	})
}

func AddToBlacklist(token string, expiration time.Duration) error {
	if RedisClient == nil {
		log.Println("错误：Redis 客户端未初始化")
		return errors.New("Redis 客户端未初始化")
	}
	key := "jti:" + token
	err := RedisClient.Set(Ctx, key, true, expiration).Err()
	if err != nil {
		log.Printf("Redis 写入失败 | Key: %s | 错误: %v\n", key, err)
		return err
	}
	// 记录成功日志：
	log.Printf("Token 已加入黑名单 | Key: %s | TTL: %v\n", key, expiration)
	return nil
}

func IsTokenRevoked(token string) bool {
	if RedisClient == nil {
		return false
	}
	exists, _ := RedisClient.Exists(Ctx, "jti:"+token).Result()
	return exists > 0
}
