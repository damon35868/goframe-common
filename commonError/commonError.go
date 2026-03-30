package commonError

import "github.com/gogf/gf/v2/errors/gcode"

const (
	DB_QUERY_ERROR_CODE   = 10000
	REDIS_LOCK_ERROR_CODE = 10001
)

var (
	DBQueryError   = gcode.New(DB_QUERY_ERROR_CODE, "数据库查询或序列化错误", nil)
	RedisLockError = gcode.New(REDIS_LOCK_ERROR_CODE, "当前操作过于频繁，请稍后重试~", nil)
)
