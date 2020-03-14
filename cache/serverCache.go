package cache

import (
	"encoding/json"
	"fmt"
	"go-im/logger"
	"go-im/model"
	redislib "go-im/redisConf"
	"go.uber.org/zap"
	"strconv"
)

const (
	serversHashKey       = "acc:hash:servers" // 全部的服务器
	serversHashCacheTime = 2 * 60 * 60        // key过期时间
	serversHashTimeout   = 3 * 60             // 超时时间
)

func getServersHashKey() (key string) {
	key = fmt.Sprintf("%s", serversHashKey)

	return
}

// 设置服务器信息
func SetServerInfo(server *model.Server, currentTime uint64) (err error) {
	key := getServersHashKey()

	value := fmt.Sprintf("%d", currentTime)

	redisClient := redislib.GetClient()
	number, err := redisClient.Do("hSet", key, server.String(), value).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do("Expire", key, serversHashCacheTime)

	return
}

// 下线服务器信息
func DelServerInfo(server *model.Server) (err error) {
	key := getServersHashKey()
	redisClient := redislib.GetClient()
	number, err := redisClient.Do("hDel", key, server.String()).Int()
	if err != nil {
		fmt.Println("DelServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	redisClient.Do("Expire", key, serversHashCacheTime)

	return
}

func GetServerAll(currentTime uint64) (servers []*model.Server, err error) {

	servers = make([]*model.Server, 0)
	key := getServersHashKey()

	redisClient := redislib.GetClient()

	val, err := redisClient.Do("hGetAll", key).Result()

	valByte, _ := json.Marshal(val)
	logger.Logger.Info("GetServerAll", zap.String("key", key), zap.String("valByte", string(valByte)))
	serverMap, err := redisClient.HGetAll(key).Result()
	if err != nil {
		fmt.Println("SetServerInfo", key, err)
		logger.Logger.Error("SetServerInfo", zap.String("key", key), zap.Error(err))
		return
	}

	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			logger.Logger.Error("GetServerAll", zap.String("key", key), zap.Error(err))
			return nil, err
		}

		// 超时
		if valueUint64+serversHashTimeout <= currentTime {
			continue
		}

		server, err := model.StringToServer(key)
		if err != nil {
			logger.Logger.Error("GetServerAll", zap.String("key", key), zap.Error(err))
			return nil, err
		}

		servers = append(servers, server)
	}

	return
}
