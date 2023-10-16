package main

import (
	"walk-server/global"
	"walk-server/router"
	"walk-server/utility"
	"walk-server/utility/initial"
	"walk-server/utility/initial/wechat"

	"github.com/gin-gonic/gin"
)

func main() {
	initial.ConfigInit() // 读取配置
	initial.DBInit()     // 初始化数据库
	initial.RedisInit()  // 初始化Redis
	initial.LimitInit()  // 初始化令牌桶
	initial.ConstantInit()
	wechat.WeChatInit()

	// 如果配置文件中开启了调试模式
	if !utility.IsDebugMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化路由
	r := initial.RouterInit()
	router.MountRoutes(r)

	// 启动服务器
	utility.StartServer(r, ":"+global.Config.GetString("server.port"))
}
