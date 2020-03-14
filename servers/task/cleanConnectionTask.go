/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-31
* Time: 15:17
 */

package task

import (
	"go-im/logger"
	"go-im/servers/websocket"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

func Init() {
	Timer(3*time.Second, 30*time.Second, cleanConnection, "", nil, nil)

}

// 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("ClearTimeoutConnections stop", zap.Any("r", r), zap.Binary("stack", debug.Stack()))
		}
	}()

	logger.Logger.Info("定时任务，清理超时连接", zap.Any("param", param))
	websocket.ClearTimeoutConnections()

	return
}
