package main

import (
	"github.com/liumkssq/webook/settings"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

func main() {
	//initViperV1()
	settings.CheckEtcdConnection()
	//initViperRemote()
	settings.InitViperV1()
	//err := logger.InitLogger()

	//fmt.Println(viper.AllKeys())
	server := InitWebServer()
	zap.L().Info("服务器启动成功")
	server.Run()
}
