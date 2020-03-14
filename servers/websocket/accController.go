package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"go-im/cache"
	"go-im/common"
	"go-im/logger"
	"go-im/model"
	"go.uber.org/zap"
	"time"
)

// ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	logger.Logger.Info("webSocket_request ping接口",
		zap.String("seq", seq),
		zap.String("Addr", client.Addr),
		zap.ByteString("message", message),
	)
	data = "pong"

	return
}

// 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &model.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		logger.Logger.Error("用户登录 解析数据失败",
			zap.String("seq", seq),
			zap.Error(err))
		return
	}

	logger.Logger.Info("webSocket_request 用户登录", zap.String("seq", seq),
		zap.String("ServiceToken", request.ServiceToken))
	// TODO::进行用户权限认证，一般是客户端传入TOKEN，然后检验TOKEN是否合法，通过TOKEN解析出来用户ID
	// 本项目只是演示，所以直接过去客户端传入的用户ID
	if request.UserId == "" || len(request.UserId) >= 20 {
		code = common.UnauthorizedUserId
		logger.Logger.Info("用户登录 非法的用户", zap.String("seq", seq),
			zap.String("UserId", request.UserId))
		return
	}

	if !InAppIds(request.AppId) {
		code = common.Unauthorized
		logger.Logger.Info("用户登录 不支持的平台", zap.String("seq", seq),
			zap.Uint32("AppId", request.AppId))
		return
	}

	if client.IsLogin() {
		fmt.Println("用户登录 用户已经登录", client.AppId, client.UserId, seq)
		logger.Logger.Info("用户登录 用户已经登录",
			zap.Uint32("AppId", request.AppId),
			zap.String("UserId", client.UserId),
			zap.String("seq", seq))
		code = common.OperationFailure

		return
	}

	client.Login(request.AppId, request.UserId, currentTime)

	// 存储数据
	userOnline := model.UserLogin(serverIp, serverPort, request.AppId, request.UserId, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		logger.Logger.Error("用户登录 SetUserOnlineInfo",
			zap.String("seq", seq), zap.Error(err))
		return
	}

	// 用户登录
	login := &login{
		AppId:  request.AppId,
		UserId: request.UserId,
		Client: client,
	}
	clientManager.Login <- login

	logger.Logger.Info("用户登录 成功",
		zap.String("Addr", client.Addr),
		zap.String("UserId", request.UserId),
		zap.String("seq", seq))
	return
}

// 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &model.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		logger.Logger.Error("心跳接口 解析数据失败",
			zap.String("seq", seq), zap.Error(err))
		return
	}

	logger.Logger.Info("webSocket_request 心跳接口",
		zap.Uint32("AppId", client.AppId),
		zap.String("UserId", client.UserId))

	if !client.IsLogin() {
		logger.Logger.Info("心跳接口 用户未登录",
			zap.Uint32("AppId", client.AppId),
			zap.String("UserId", client.UserId),
			zap.String("seq", seq))
		code = common.NotLoggedIn

		return
	}

	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			logger.Logger.Info("心跳接口 用户未登录",
				zap.Uint32("AppId", client.AppId),
				zap.String("UserId", client.UserId),
				zap.String("seq", seq))
			return
		} else {
			code = common.ServerError
			logger.Logger.Error("心跳接口 GetUserOnlineInfo",
				zap.String("seq", seq),
				zap.Uint32("AppId", client.AppId),
				zap.String("UserId", client.UserId),
				zap.Error(err))
			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		logger.Logger.Error("心跳接口 SetUserOnlineInfo",
			zap.String("seq", seq),
			zap.Uint32("AppId", client.AppId),
			zap.String("UserId", client.UserId),
			zap.Error(err))
		return
	}

	return
}
