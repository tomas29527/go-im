/**
* Created by GoLand.
* User: link1st
* Date: 2019-07-25
* Time: 12:11
 */

package home

import (
	"github.com/gin-gonic/gin"
	"go-im/common"
	"go-im/config"
	controllers "go-im/controller"
	"net/http"
)

// 查看用户是否在线
func Index(c *gin.Context) {
	appId := c.Query("appId")

	data := gin.H{
		"title":        "聊天首页",
		"httpUrl":      config.GlobalConfig.App.HttpUrl,
		"webSocketUrl": config.GlobalConfig.App.WebSocketUrl,
		"appId":        appId,
	}
	c.HTML(http.StatusOK, "index.tpl", data)
}

func Test(c *gin.Context) {
	data := make(map[string]interface{})
	data["aa"] = "ccc"
	controllers.Response(c, common.OK, "", data)
}
