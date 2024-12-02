package ioc

import (
	"fmt"
	"github.com/liumkssq/webook/internal/repository/dao"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	//var cfg = Config{
	//	DSN: "root:root@tcp(localhost:13316)/webook_default",
	//}
	var cfg Config
	err := viper.UnmarshalKey("db", &cfg)
	fmt.Println(cfg.DSN)
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		//Logger: glogger.
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
