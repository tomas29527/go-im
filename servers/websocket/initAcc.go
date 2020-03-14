/**
 * Created by GoLand.
 * User: link1st
 * Date: 2019-07-25
 * Time: 16:04
 */

package websocket

import (
	"github.com/gorilla/websocket"
	"go-im/config"
	"go-im/logger"
	"go-im/model"
	"go-im/util"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var (
	clientManager = NewClientManager() // 管理者
	appIds        = []uint32{101, 102} // 全部的平台

	serverIp   string
	serverPort string
)

func GetAppIds() []uint32 {

	return appIds
}

func GetServer() (server *model.Server) {
	server = model.NewServer(serverIp, serverPort)

	return
}

func IsLocal(server *model.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}

	return
}

func InAppIds(appId uint32) (inAppId bool) {

	for _, value := range appIds {
		if value == appId {
			inAppId = true

			return
		}
	}

	return
}

// 启动程序
func StartWebSocket() {
	serverIp = util.GetServerIp()

	webSocketPort := config.GlobalConfig.App.WebSocketPort
	//rpcPort := viper.GetString("app.rpcPort")

	//serverPort = rpcPort

	http.HandleFunc("/acc", wsPage)

	// 添加处理程序
	go clientManager.start()

	logger.Logger.Info("WebSocket 启动程序成功",
		zap.String("serverIp", serverIp),
		zap.String("webSocketPort", webSocketPort))

	http.ListenAndServe(":"+webSocketPort, nil)
}

func wsPage(w http.ResponseWriter, req *http.Request) {

	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		logger.Logger.Info("升级协议", zap.Any("ua:",
			r.Header["User-Agent"]),
			zap.Any("referer", r.Header["Referer"]))
		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	logger.Logger.Info("webSocket 建立连接:", zap.String("RemoteAddr", conn.RemoteAddr().String()))
	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
