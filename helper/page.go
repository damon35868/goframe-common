package helper

import (
	"github.com/damon35868/goframe-common/dto"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func PageBuilder[T any](model *gdb.Model, page, pageSize int, conditions ...func(model *gdb.Model) *gdb.Model) (res *dto.PageResDto[T], err error) {
	res = &dto.PageResDto[T]{
		Items: []T{},
	}
	if len(conditions) > 0 {
		model = conditions[0](model)
	}
	if err := model.Page(page, pageSize).ScanAndCount(&res.Items, &res.TotalCount, false); err != nil {
		return nil, gerror.NewCode(gcode.New(10000, "数据库查询或序列化错误", nil))
	}

	res.HasNextPage = HasNextPage(page, pageSize, res.TotalCount)
	return res, nil
}

func HasNextPage(page, pageSize, totalCount int) bool {
	return page*pageSize < totalCount
}
