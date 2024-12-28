package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

// RedisDistributedLock 是一个简单的分布式锁结构
type RedisDistributedLock struct {
	client  *redis.Client
	lockKey string
	value   string // 锁的唯一标识符，通常是客户端ID或UUID
	ttl     time.Duration
}

// NewRedisDistributedLock 创建一个新的分布式锁实例
func NewRedisDistributedLock(client *redis.Client, lockKey, value string, ttl time.Duration) *RedisDistributedLock {
	return &RedisDistributedLock{
		client:  client,
		lockKey: lockKey,
		value:   value,
		ttl:     ttl,
	}
}

// TryLock 使用Lua脚本尝试获取锁，成功返回true，失败返回false
func (r *RedisDistributedLock) TryLock() bool {
	// Lua脚本：只有当锁不存在时才设置它，并设置过期时间
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == false then
			redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
			return 1
		else
			return 0
		end`)

	result, err := script.Run(ctx, r.client, []string{r.lockKey}, r.value, fmt.Sprintf("%d", int(r.ttl.Seconds()))).Result()
	if err != nil {
		fmt.Printf("Failed to acquire lock with Lua script: %v", err)
		return false
	}

	return result.(int64) == 1
}

// Unlock 释放锁
func (r *RedisDistributedLock) Unlock() error {
	// 使用Lua脚本保证原子性：只有当锁的值匹配时才删除
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end`)
	err := script.Run(ctx, r.client, []string{r.lockKey}, r.value).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		fmt.Printf("Failed to release lock: %v", err)
		return err
	}
	return nil
}
