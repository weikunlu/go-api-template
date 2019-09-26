package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/weikunlu/go-api-template/config"
	"sync"
	"time"
)

// Singleton setup
var singleton *cacheService
var once sync.Once

type cacheService struct {
	cacheCli *redis.Client
}

func GetCacheServiceInstance() *cacheService {
	once.Do(func() {
		cfg := config.GetAppConfig()

		fmt.Printf("init redis connection\n")
		redisCli := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDb,
		})

		pong, err := redisCli.Ping().Result()
		if err != nil {
			panic(fmt.Sprintf("fail to connect redis %v", err.Error()))
		}
		fmt.Printf("redis PING %v\n", pong)

		singleton = &cacheService{
			cacheCli: redisCli,
		}
	})
	return singleton
}

func (s *cacheService) PushTaskEvent(taskSid string, event string) (err error) {
	key := fmt.Sprintf("twilio:tasks:%s", taskSid)
	err = s.cacheCli.LPush(key, event).Err()
	if err != nil {
		return
	}
	_, err = s.cacheCli.Expire(key, time.Hour*24).Result()
	return
}

func (s *cacheService) GetAllTaskEvent(taskSid string) (results []string, err error) {
	key := fmt.Sprintf("twilio:tasks:%s", taskSid)
	results, err = s.cacheCli.LRange(key, 0, -1).Result()
	return
}

func (s *cacheService) DelTaskEvent(taskSid string) (num int64, err error) {
	key := fmt.Sprintf("twilio:tasks:%s", taskSid)
	num, err = s.cacheCli.Del(key).Result()
	return
}
