package main

import (
	"fmt"
	"github.com/temprory/redis"
	"github.com/temprory/util"
	"time"
)

var (
	CONF_BAC_SVR_1 = ""

	//发布订阅key，用来通知程序更新
	CONF_PUB_KEY = "CONF_PUB_KEY"
	//配置存储的key fieled，用于读取配置
	CONF_UPDATE_KEY              = "CONF_UPDATE_KEY"
	CONF_UPDATE_FIELED_BAC_SVR_1 = "CONF_UPDATE_FIELED_BAC_SVR_1"
)

func onUpdateBacSvr1(tag string, value string) error {
	if value == "" {
		fmt.Printf("--- onUpdateBacSvr1 failed: '%v', '%v'\n", tag, value)
		return fmt.Errorf("invalid config string")
	}

	CONF_BAC_SVR_1 = value

	fmt.Printf("--- onUpdateBacSvr1 success: '%v', '%v'\n", tag, value)

	return nil
}

func main() {

	rds := redis.NewRedis(redis.RedisConf{
		Addr:     "127.0.0.1:6379",
		PoolSize: 10,
	})

	loader := util.NewConfigLoader(rds.Client(), CONF_PUB_KEY, time.Second*10)
	loader.Add(CONF_UPDATE_KEY, CONF_UPDATE_FIELED_BAC_SVR_1, onUpdateBacSvr1)

	i := 0
	for {
		i++
		time.Sleep(time.Second * 5)
		//设置配置值
		confStr := fmt.Sprintf("[config %v]", i)
		rds.Client().HSet(CONF_UPDATE_KEY, CONF_UPDATE_FIELED_BAC_SVR_1, confStr)
		//通知更新
		rds.Client().Publish(CONF_PUB_KEY, CONF_UPDATE_KEY+":"+CONF_UPDATE_FIELED_BAC_SVR_1)
	}
}
