package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"go-im/logger"
	"go-im/model"
	redislib "go-im/redisConf"
	"go.uber.org/zap"
)

const (
	userOnlinePrefix    = "acc:user:online:" // 用户在线状态
	userOnlineCacheTime = 24 * 60 * 60
)

/*********************  查询用户是否在线  ************************/
func getUserOnlineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", userOnlinePrefix, userKey)

	return
}

func GetUserOnlineInfo(userKey string) (userOnline *model.UserOnline, err error) {
	redisClient := redislib.GetClient()

	key := getUserOnlineKey(userKey)

	data, err := redisClient.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			logger.Logger.Error("GetUserOnlineInfo",
				zap.String("userKey", userKey),
				zap.Error(err))
			return
		}

		logger.Logger.Error("GetUserOnlineInfo",
			zap.String("userKey", userKey),
			zap.Error(err))
		return
	}

	userOnline = &model.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		logger.Logger.Error("获取用户在线数据 json Unmarshal",
			zap.String("userKey", userKey),
			zap.Error(err))
		return
	}

	logger.Logger.Info("获取用户在线数据",
		zap.String("userKey", userKey),
		zap.Uint64("LoginTime", userOnline.LoginTime),
		zap.Uint64("HeartbeatTime", userOnline.HeartbeatTime),
		zap.String("AccIp", userOnline.AccIp),
		zap.Bool("IsLogoff", userOnline.IsLogoff))
	return
}

// 设置用户在线数据
func SetUserOnlineInfo(userKey string, userOnline *model.UserOnline) (err error) {
	redisClient := redislib.GetClient()
	key := getUserOnlineKey(userKey)
	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		logger.Logger.Error("设置用户在线数据 json Marshal", zap.String("key", key), zap.Error(err))
		return
	}

	_, err = redisClient.Do("setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		logger.Logger.Error("设置用户在线数据", zap.String("key", key), zap.Error(err))
		return
	}

	return
}
