package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type CacheService struct {
	rdb *redis.Client
	sf  singleflight.Group
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	cs := NewCacheService(rdb)
	data, err := cs.GetUserData(context.Background(), "123")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(data)
}

func NewCacheService(rdb *redis.Client) *CacheService {
	return &CacheService{
		rdb: rdb,
		sf:  singleflight.Group{},
	}
}

func (s *CacheService) GetUserData(ctx context.Context, userID string) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("user:%s", userID)

	// 先尝试从缓存获取
	userData, err := s.rdb.HGetAll(ctx, cacheKey).Result()
	if err == nil && len(userData) > 0 {
		// 缓存命中，直接返回
		result := make(map[string]interface{})
		for k, v := range userData {
			result[k] = v
		}
		return result, nil
	}

	// 缓存未命中，使用 singleflight 合并并发请求
	v, err, _ := s.sf.Do(userID, func() (interface{}, error) {
		// 再次检查缓存（双重检查避免竞态条件）
		userData, err := s.rdb.HGetAll(ctx, cacheKey).Result()
		if err == nil && len(userData) > 0 {
			result := make(map[string]interface{})
			for k, v := range userData {
				result[k] = v
			}
			return result, nil
		}

		// 从数据库加载数据（模拟耗时操作）
		data, err := loadUserFromDB(ctx, userID)
		if err != nil {
			return nil, err
		}

		// 写入缓存
		cacheData := make(map[string]interface{})
		for k, v := range data {
			cacheData[k] = fmt.Sprintf("%v", v)
		}
		s.rdb.HSet(ctx, cacheKey, cacheData)
		s.rdb.Expire(ctx, cacheKey, time.Minute*15) // 设置过期时间

		return data, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(map[string]interface{}), nil
}

// 模拟从数据库加载用户信息
func loadUserFromDB(ctx context.Context, userID string) (map[string]interface{}, error) {
	// 模拟数据库访问延迟
	time.Sleep(time.Millisecond * 200)

	// 模拟数据库返回结果
	return map[string]interface{}{
		"id":    userID,
		"name":  "用户" + userID,
		"email": "user" + userID + "@example.com",
	}, nil
}
