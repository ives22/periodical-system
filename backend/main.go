package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"periodical/conf"
	periodicalApi "periodical/internal/apps/periodical/api"
	periodicalImpl "periodical/internal/apps/periodical/impl"
	"periodical/internal/apps/token/api"
	tokenImpl "periodical/internal/apps/token/impl"
	userApi "periodical/internal/apps/user/api"
	//userImpl "promalert/internal/apps/user/impl"
	userImpl "periodical/internal/apps/user/impl"
)

func main() {
	if err := conf.LoadConfigFromYaml("etc/config.yaml"); err != nil {
		panic(err)
	}

	// 初始化控制器
	// user controller
	userSvcImpl := userImpl.NewUserImpl()
	// token controller
	tokenSvcImpl := tokenImpl.NewTokenImpl(userSvcImpl)

	// periodical controller
	periodicalSvc := periodicalImpl.NewPeriodicalImpl()

	////
	//req := user.NewCreateUserRequest()
	//req.Username = "admin"
	//req.Password = "123456"
	//
	//u, err := userSvcImpl.CreateUser(context.Background(), req)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(u)

	// 实例化期刊API对象
	// 期刊 api handler
	pApi := periodicalApi.NewPeriodicalApiHandler(periodicalSvc)
	// token api handler
	tkApi := api.NewTokenApiHandler(tokenSvcImpl)
	// User api handler
	uApi := userApi.NewUserApiHandler(userSvcImpl)

	// 启动http协议服务器，注册handler路由
	r := gin.Default()

	r.Use(cors.Default())

	tkApi.Registry(r) // token handler
	pApi.Registry(r)  // 期刊handler
	uApi.Registry(r)  // 用户handler

	r.Run(conf.C().App.HttpAddr())
}
