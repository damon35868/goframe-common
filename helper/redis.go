package helper

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type RedisLock struct {
	Ctx        context.Context // context
	Key        string          // 锁的 Key
	Value      string          // 锁的值,通常为一个唯一的标识,如 UUID
	ExpireTime int64           // 锁的过期时间(秒)
}

// Lock 获取锁
// 返回 true: 获取锁成功, false: 获取锁失败
func (lock *RedisLock) Lock() bool {
	ok, err := g.Redis().SetNX(lock.Ctx, lock.Key, lock.Value)
	if err != nil {
		g.Log().Error(lock.Ctx, "SetNX error", err)
		return false
	}
	if !ok {
		// 锁已经存在，获取失败
		return false
	}
	// 设置成功后设置过期时间，Expire 参数需要 int64
	_, err = g.Redis().Expire(lock.Ctx, lock.Key, lock.ExpireTime)
	if err != nil {
		g.Log().Error(lock.Ctx, "Expire error", err)
		// 如果设置过期失败，删除锁
		_, _ = g.Redis().Del(lock.Ctx, lock.Key)
		return false
	}
	return true
}

// Unlock 释放锁
// 使用 Lua 脚本保证只有持有者才能释放锁
func (lock *RedisLock) Unlock() bool {
	// 使用 Lua 脚本释放锁
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
	`
	args := []interface{}{lock.Value}
	result, err := g.Redis().Eval(lock.Ctx, script, 1, []string{lock.Key}, args)
	if err != nil {
		g.Log().Error(lock.Ctx, "释放锁失败", err)
		return false
	}
	return result.Int() == 1
}
