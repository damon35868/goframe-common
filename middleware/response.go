package middleware

import (
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err             = r.GetError()
		res             = r.GetHandlerResponse()
		code gcode.Code = gcode.CodeOK
	)

	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		// 错误情况下，如果 res 为 nil，只输出 code 和 message，不包含 data
		if res == nil || (reflect.ValueOf(res).Kind() == reflect.Ptr && reflect.ValueOf(res).IsNil()) {
			r.Response.WriteJson(g.Map{
				"code":    code.Code(),
				"message": err.Error(),
			})
		}
		return
	}
	// 404处理
	if r.Response.Status == 404 {
		r.Response.WriteJson(g.Map{
			"code":    http.StatusNotFound,
			"message": http.StatusText(http.StatusNotFound),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"code":    code.Code(),
		"message": gcode.CodeOK.Message(),
		"data":    res,
	})
}
