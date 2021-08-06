
package main

import (
	"GF_PROJECT_NAME/core"
	"GF_PROJECT_NAME/global"
	"GF_PROJECT_NAME/initialize"
	"GF_PROJECT_NAME/model"
	"GF_PROJECT_NAME/router"
)

// @title Swagger Example API
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /
func main() {
	global.VP = core.Viper()      // 初始化Viper
	global.LOG = core.Zap()       // 初始化zap日志库
	global.DB = initialize.Gorm() // gorm连接数据库
//	initialize.Timer()
	if global.DB != nil {
		models := make([]interface{}, 0)
		models = append(models, model.System{})
		initialize.MysqlTables(global.DB, models) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.DB.DB()
		defer db.Close()
	}
	wholeRouter := initialize.Routers()

        initFuncs := make([]initialize.InitRouterGroupFunc, 0)
        initFuncs = append(initFuncs, router.InitSystemRouter)
	initialize.AddRouterGroup(wholeRouter, initFuncs)

	core.RunWindowsServer(wholeRouter)
}
