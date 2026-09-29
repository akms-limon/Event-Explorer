package middleware

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

func RequestID(ctx *context.Context) {
	requestID := fmt.Sprintf("%d", time.Now().UnixNano())

	ctx.Input.SetData("request_id", requestID)
}