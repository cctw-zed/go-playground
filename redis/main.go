package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// 创建 Redis 客户端
func createRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 无密码
		DB:       0,  // 默认 DB
	})
	return client
}

// String 类型操作示例
func stringOperationExample(ctx context.Context, rdb *redis.Client) error {
	// 1. 基本的 Set 和 Get 操作
	err := rdb.Set(ctx, "user:1:name", "张三", 0).Err()
	if err != nil {
		return fmt.Errorf("set key failed: %v", err)
	}

	// 2. 使用 SetNX 实现互斥锁
	locked, err := rdb.SetNX(ctx, "lock:order:12345", "1", time.Second*30).Result()
	if err != nil {
		return fmt.Errorf("acquire lock failed: %v", err)
	}
	if !locked {
		return fmt.Errorf("failed to acquire lock, someone else holds it")
	}

	// 3. 批量操作
	pipe := rdb.Pipeline()
	pipe.Set(ctx, "user:1:age", "25", 0)
	pipe.Set(ctx, "user:1:city", "北京", 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline execute failed: %v", err)
	}

	// 4. 带过期时间的缓存操作
	err = rdb.Set(ctx, "verification:code:13812345678", "123456", time.Minute*5).Err()
	if err != nil {
		return fmt.Errorf("set verification code failed: %v", err)
	}

	return nil
}

// 实现分布式锁
func distributedLock(ctx context.Context, rdb *redis.Client, lockKey string) (bool, func()) {
	// 获取锁，设置过期时间为 30 秒
	success, err := rdb.SetNX(ctx, lockKey, "1", time.Second*30).Result()
	if err != nil || !success {
		return false, nil
	}

	// 返回解锁函数
	unlock := func() {
		rdb.Del(ctx, lockKey)
	}

	return true, unlock
}

func main() {
	ctx := context.Background()
	rdb := createRedisClient()
	defer rdb.Close()

	//err := stringOperationExample(ctx, rdb)
	//if err != nil {
	//	fmt.Println(err)
	//}
	// 使用分布式锁的示例
	locked, unlock := distributedLock(ctx, rdb, "lock:order:12345")
	if locked {
		defer unlock()
		// 执行需要加锁的业务逻辑
		fmt.Println("获取锁成功，执行业务逻辑")
	} else {
		fmt.Println("获取锁失败")
	}
}
