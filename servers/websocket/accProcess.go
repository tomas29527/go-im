/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-27
 * Time: 14:38
 */

package websocket

import (
	"encoding/json"
	"go-im/common"
	"go-im/logger"
	"go-im/model"
	"go.uber.org/zap"
	"sync"
)

type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// 注册
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}

// 处理数据
func ProcessData(client *Client, message []byte) {
	logger.Logger.Info("处理数据", zap.String("Addr", client.Addr), zap.Binary("message", message))
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("处理数据 stop", zap.Any("r", r))
		}
	}()

	request := &model.Request{}

	err := json.Unmarshal(message, request)
	if err != nil {
		logger.Logger.Error("处理数据 json Unmarshal", zap.Error(err))
		client.SendMsg([]byte("数据不合法"))
		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		logger.Logger.Error("处理数据 json Marshal", zap.Error(err))
		client.SendMsg([]byte("处理数据失败"))
		return
	}

	seq := request.Seq
	cmd := request.Cmd

	var (
		code uint32
		msg  string
		data interface{}
	)

	// request
	logger.Logger.Info("acc_request", zap.String("cmd", cmd), zap.String("Addr", client.Addr))
	// 采用 map 注册的方式
	if value, ok := getHandlers(cmd); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = common.RoutingNotExist
		logger.Logger.Info("处理数据 路由不存在", zap.String("Addr", client.Addr), zap.String("cmd", cmd))
	}

	msg = common.GetErrorMessage(code, msg)

	responseHead := model.NewResponseHead(seq, cmd, code, msg, data)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		logger.Logger.Error("处理数据 json Marshal", zap.Error(err))
		return
	}

	client.SendMsg(headByte)
	logger.Logger.Info("acc_response send",
		zap.String("Addr", client.Addr),
		zap.Uint32("AppId", client.AppId),
		zap.String("UserId", client.UserId),
		zap.String("cmd", cmd),
		zap.Uint32("code", code))
	return
}
