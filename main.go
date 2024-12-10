package main

import (
	"github.com/liumkssq/webook/settings"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

func main() {
	//initViperV1()
	//settings.CheckEtcdConnection()
	//initViperRemote()
	settings.InitViperV1()
	//err := logger.InitLogger()

	//fmt.Println(viper.AllKeys())
	app := InitWebServer()
	for _,c := app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	zap.L().Info("服务器启动成功")
	app.Run()
}
