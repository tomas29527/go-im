package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-im/config"
	"go-im/logger"
	redislib "go-im/redisConf"
	"go-im/routers"
	"go-im/servers/task"
	"go-im/servers/websocket"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	//初始化配置
	err := config.Initconfig()
	if err != nil {
		fmt.Println("=======服务启动失败,配置加载失败======", err)
		return
	}
	//初始化日志
	logger.InitLog(config.GlobalConfig.App.LogFile)
	//初始化redis
	initRedis()

	logger.Logger.Info("===日志==")

	//初始化gin服务
	router := gin.Default()
	// 初始化路由
	routers.Init(router)
	routers.WebsocketInit()

	// 定时任务
	task.Init()
	// 服务注册
	task.ServerInit()

	go websocket.StartWebSocket()

	go open()

	httpPort := config.GlobalConfig.App.HttpPort
	http.ListenAndServe(":"+httpPort, router)
}

func initRedis() {
	redislib.ExampleNewClient()
}

func open() {

	time.Sleep(1000 * time.Millisecond)

	httpUrl := config.GlobalConfig.App.HttpUrl
	httpUrl = "http://" + httpUrl + "/home/index"

	fmt.Println("访问页面体验:", httpUrl)

	cmd := exec.Command("open", httpUrl)
	cmd.Output()
}
