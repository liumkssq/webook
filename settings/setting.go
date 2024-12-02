package settings

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
)

func InitViperV1() {
	cfile := pflag.String("config",
		"config/config.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉你差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	//viper.SetDefault("db.mysql.dsn",
	//	"root:root@tcp(localhost:3306)/mysql")
	//viper.SetConfigFile("config/dev.yaml")
	//viper.KeyDelimiter("-")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func InitViper() {
	viper.SetDefault("db.mysql.dsn",
		"root:root@tcp(localhost:3306)/mysql")
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go, .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是 yaml 格式
	// 现实中，有很多格式，JSON，XML，YAML，TOML，ini
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./config")
	//viper.AddConfigPath("/tmp/config")
	//viper.AddConfigPath("/etc/webook")
	// 读取配置到 viper 里面，或者你可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//otherViper := viper.New()
	//otherViper.SetConfigName("myjson")
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigType("json")
}

func CheckEtcdConnection() {
	resp, err := http.Get("http://127.0.0.1:12379/health")
	if err != nil {
		fmt.Println("Error connecting to Etcd:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Connected to Etcd successfully")
	} else {
		fmt.Println("Failed to connect to Etcd, status code:", resp.StatusCode)
	}
}

func InitViperRemote() {
	err := viper.AddRemoteProvider("etcd3",
		// 通过 webook 和其他使用 etcd 的区别出来
		"http://localhost:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func InitViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13316)/webook"

redis:
  addr: "localhost:6379"
`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}
