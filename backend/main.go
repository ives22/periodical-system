package main

import (
	"github.com/gin-gonic/gin"

	"periodical/conf"
	aiApi "periodical/internal/apps/ai/api"
	aiImpl "periodical/internal/apps/ai/impl"
	periodicalApi "periodical/internal/apps/periodical/api"
	periodicalImpl "periodical/internal/apps/periodical/impl"
	roleApi "periodical/internal/apps/role/api"
	roleImpl "periodical/internal/apps/role/impl"
	"periodical/internal/apps/token/api"
	tokenImpl "periodical/internal/apps/token/impl"
	userApi "periodical/internal/apps/user/api"
	userImpl "periodical/internal/apps/user/impl"

	"github.com/gin-contrib/cors"
)

func main() {
	if err := conf.LoadConfigFromYaml("etc/config.yaml"); err != nil {
		panic(err)
	}

	// 初始化控制器
	// role & permission controller (含自动建表与默认RBAC数据初始化)
	roleSvcImpl := roleImpl.NewRoleImpl()
	// user controller
	userSvcImpl := userImpl.NewUserImpl(roleSvcImpl)
	// token controller
	tokenSvcImpl := tokenImpl.NewTokenImpl(userSvcImpl)
	// periodical controller
	periodicalSvc := periodicalImpl.NewPeriodicalImpl()
	// ai agent controller (含大模型渠道配置、会话管理与Function Calling编排)
	aiSvcImpl := aiImpl.NewAIImpl(periodicalSvc)

	// 实例化各模块API Handler
	pApi := periodicalApi.NewPeriodicalApiHandler(periodicalSvc)
	tkApi := api.NewTokenApiHandler(tokenSvcImpl)
	uApi := userApi.NewUserApiHandler(userSvcImpl, roleSvcImpl)
	rApi := roleApi.NewRoleApiHandler(roleSvcImpl)
	aApi := aiApi.NewAIApiHandler(aiSvcImpl)

	// 启动http协议服务器，注册handler路由
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	tkApi.Registry(r)              // 认证 handler
	pApi.Registry(r, tokenSvcImpl) // 期刊 handler
	uApi.Registry(r, tokenSvcImpl) // 用户 handler
	rApi.Registry(r, tokenSvcImpl) // 角色与权限 handler
	aApi.Registry(r, tokenSvcImpl) // AI Agent 与大模型管理 handler

	r.Run(conf.C().App.HttpAddr())
}
