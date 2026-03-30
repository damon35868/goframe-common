package helper

import (
	"context"
	"time"

	"github.com/damon35868/goframe-common/cache"
	"github.com/damon35868/goframe-common/commonError"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type RedisLock struct {
	Ctx          context.Context // context
	Key          string          // 锁的 Key
	Value        string          // 锁的值,通常为一个唯一的标识,如 UUID
	ExpireTime   int64           // 锁的过期时间(秒)
	WaitTimeout  int64           // 等待锁释放的最长时间(秒),0 表示不等待,获取失败立即返回
	WaitInterval time.Duration   // 轮询获取锁的间隔时间,默认 100ms
	ErrMsg       string          // 锁失败的错误信息
}

// Lock 获取锁
// 返回 true: 获取锁成功, false: 获取锁失败
func (lock *RedisLock) Lock() bool {
	// 如果没有设置等待间隔，使用默认 100ms
	if lock.WaitInterval <= 0 {
		lock.WaitInterval = 100 * time.Millisecond
	}

	var deadline time.Time
	if lock.WaitTimeout > 0 {
		deadline = time.Now().Add(time.Duration(lock.WaitTimeout) * time.Second)
	}

	for {
		ok, err := g.Redis().SetNX(lock.Ctx, lock.Key, lock.Value)
		if err != nil {
			g.Log().Error(lock.Ctx, "SetNX error", err)
			return false
		}
		if ok {
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
		// 锁已经存在
		if lock.WaitTimeout <= 0 {
			// 不等待，立即返回失败
			return false
		}
		// 检查是否超时
		if time.Now().After(deadline) {
			// 等待超时，返回失败
			return false
		}
		// 等待一段时间后重试
		time.Sleep(lock.WaitInterval)
	}
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

func (lock *RedisLock) TryLock(action func(), releaseLock ...bool) error {
	if lock.Lock() {
		unLock := true
		if len(releaseLock) > 0 {
			unLock = releaseLock[0]
		}
		if unLock {
			defer lock.Unlock()
		}
		action()
	} else {
		if lock.ErrMsg != "" {
			return gerror.NewCode(gcode.New(commonError.REDIS_LOCK_ERROR_CODE, lock.ErrMsg, nil))
		} else {
			return gerror.NewCode(commonError.RedisLockError)
		}
	}
	return nil
}

func CacheRemember[T any](ctx context.Context, key string, duration time.Duration, action func() *T) (res *T, err error) {
	rel, err := cache.Redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if !rel.IsNil() {
		err = rel.Struct(&res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	value := action()
	if err := cache.Redis.Set(ctx, key, value, duration); err != nil {
		return nil, err
	}

	return value, nil
}
