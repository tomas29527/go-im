package systems

import (
	"github.com/gin-gonic/gin"
	"go-im/common"
	controllers "go-im/controller"
	"go-im/logger"
	"go-im/servers/websocket"
	"go.uber.org/zap"
	"runtime"
)

// 查询系统状态
func Status(c *gin.Context) {

	isDebug := c.Query("isDebug")
	logger.Logger.Info("http_request 查询系统状态", zap.String("isDebug", isDebug))
	data := make(map[string]interface{})

	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// goroutine 的数量
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU

	// ClientManager 信息
	data["managerInfo"] = websocket.GetManagerInfo(isDebug)

	controllers.Response(c, common.OK, "", data)
}
