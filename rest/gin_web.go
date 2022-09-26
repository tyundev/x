package rest

import (
	"runtime/debug"

	//"github.com/reiwav/x/rest/validator"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

const STATUS_OK = 200

type JsonRender struct {
}

func (r *JsonRender) SendData(ctx *gin.Context, data interface{}) {
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   data,
		"status": "success",
		"code":   200,
	})
}

func (r *JsonRender) SendString(ctx *gin.Context, data interface{}) {
	ctx.JSON(STATUS_OK, data)
}

func (r *JsonRender) SendDataNotFound(ctx *gin.Context, data interface{}, isNotFound bool) {
	var status = "success"
	if isNotFound {
		status = "error"
	}
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   data,
		"status": status,
		"code":   200,
	})
}

func (r *JsonRender) SendDataError(ctx *gin.Context, data interface{}, err error) {
	var isCheck = AssertNexNotFound(err)
	var status = "success"
	if isCheck {
		status = "error"
	}
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   data,
		"status": status,
		"code":   200,
	})
}

func (r *JsonRender) SendError(ctx *gin.Context, err error) {
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"message": err.Error(),
		"status":  "error",
		"code":    200,
	})
}

func (r *JsonRender) DecodeBody(ctx *gin.Context, data interface{}) {
	AssertNil(BadRequest(ctx.BindJSON(&data).Error()))
	//AssertNil(BadRequest(validator.Validate(data).Error()))

}
func (r *JsonRender) Success(ctx *gin.Context) {
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   nil,
		"status": "success",
		"code":   200,
	})
}

func Recover() {
	if r := recover(); r != nil {
		glog.Error(r, string(debug.Stack()))
	}
}
