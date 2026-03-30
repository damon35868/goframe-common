package cache

import (
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

var Redis = gcache.New()

func init() {
	redis, err := gredis.New(&gredis.Config{
		Address: g.Cfg().MustGet(gctx.GetInitCtx(), "redis.cache.address").String(),
		Db:      g.Cfg().MustGet(gctx.GetInitCtx(), "redis.cache.db").Int(),
	})
	if err != nil {
		panic(err)
	}
	Redis.SetAdapter(gcache.NewAdapterRedis(redis))
}
