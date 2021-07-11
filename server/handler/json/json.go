package json

import (
	"net/http"

	"github.com/sinksmell/toy-go/protocol"
	converter "github.com/sinksmell/toy-go/tool/converter"

	"github.com/douyu/jupiter/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ConvertToGo(ctx echo.Context) error {
	var (
		err  error
		req  = &protocol.JsonConvertReq{}
		resp = &protocol.JsonConvertResp{}
	)

	if err = ctx.Bind(req); err != nil {
		xlog.Error("bind convert json str err", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest,
			&protocol.JsonConvertError{Message: err.Error()})
	}

	err = jsoniter.Unmarshal([]byte(req.JsonStr), &map[string]interface{}{})
	if err != nil {
		xlog.Error("check json str err", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest,
			&protocol.JsonConvertError{Message: err.Error()})
	}

	resp.GoStructStr = converter.NewJsonConverter().GenStruct(req.JsonStr)

	return ctx.JSON(http.StatusOK, resp)
}
