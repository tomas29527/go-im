/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-03
* Time: 15:44
 */

package task

import (
	"go-im/cache"
	"go-im/logger"
	"go-im/servers/websocket"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

func ServerInit() {
	Timer(2*time.Second, 60*time.Second, server, "", serverDefer, "")
}

// 服务注册
func server(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("服务注册 stop", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
		}
	}()

	server := websocket.GetServer()
	currentTime := uint64(time.Now().Unix())
	logger.Logger.Info("定时任务，服务注册",
		zap.Any("param", param),
		zap.Any("server", server),
		zap.Uint64("currentTime", currentTime))
	cache.SetServerInfo(server, currentTime)

	return
}

// 服务下线
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("服务下线 stop", zap.Any("r", r), zap.String("stack", string(debug.Stack())))
		}
	}()

	logger.Logger.Info("服务下线", zap.Any("param", param))

	server := websocket.GetServer()
	cache.DelServerInfo(server)

	return
}
