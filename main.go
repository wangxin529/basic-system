package main

import (
	"elevate-hub/config"
	"elevate-hub/router"
	"elevate-hub/svc"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func main() {

	//fmt.Println(runtime.Version())
	// 1. 初始化配置
	var conf config.Config
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed, err:%v\n", err)
	}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("config init failed, err:%v\n", err)
	}
	// 2. 配置上下文
	ctx := svc.NewServiceContext(conf)
	//3.生成 gin服务
	engine := gin.Default()
	//4. 生成路由
	router.Router(engine, ctx)
	// 5. 启动服务
	err = engine.Run(fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatalf("http server run failed, err:%v\n", err)
	}
}
