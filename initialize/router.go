package initialize

import (
//	_ "GF_PROJECT_NAME/docs"
	"GF_PROJECT_NAME/global"
	// "GF_PROJECT_NAME/middleware"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// 初始化总路由

func Routers() *gin.Engine {
	var Router = gin.Default()
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了
	global.LOG.Info("use middleware logger")
	// 跨域
	//Router.Use(middleware.Cors()) // 如需跨域可以打开
	global.LOG.Info("use middleware cors")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.LOG.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用
	global.LOG.Info("router register success")
	return Router
}

type InitRouterGroupFunc func(Router *gin.RouterGroup)

func AddRouterGroup(router *gin.Engine, initFuncs []InitRouterGroupFunc) {
	PrivateGroup := router.Group("")
	// PrivateGroup.Use(middleware.JWTAuth())
	{
		for _, initFunc := range initFuncs {
			initFunc(PrivateGroup)	
		}
	}
}
