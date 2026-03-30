package commonError

import "github.com/gogf/gf/v2/errors/gcode"

const (
	DB_QUERY_ERROR_CODE = 10000
)

var (
	DBQueryError = gcode.New(DB_QUERY_ERROR_CODE, "数据库查询或序列化错误", nil)
)
