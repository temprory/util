package util

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/kataras/golog"
)

type ConfigLoader struct {
	sync.Mutex
	redisCli       *redis.Client
	pubsubKey      string
	updateInterval time.Duration
	updateTasks    map[string]func()
}

func (loader *ConfigLoader) Add(configKey string, configFieled string, onUpdate func(string, string)) {
	loader.Lock()
	defer loader.Unlock()

	timer := time.NewTimer(loader.updateInterval)

	update := func() {
		Safe(func() {
			confStr, err := loader.redisCli.HGet(configKey, configFieled).Result()
			if err != nil {
				golog.Infof("ConfigLoader load config %v failed: %v", configKey, err)
				return
			}
			if len(confStr) > 0 {
				onUpdate(configKey+":"+configFieled, confStr)
			}
		})

		timer.Reset(loader.updateInterval)
	}

	loader.updateTasks[configKey+":"+configFieled] = update

	Go(func() {
		for {
			select {
			case _, ok := <-timer.C:
				if !ok {
					return
				}
				update()
			}
		}
	})

	golog.Infof("ConfigLoader, Add Item, configKey: %v, configFieled: %v", configKey, configFieled)

	update()
}

func NewConfigLoader(redisCli *redis.Client, pubsubKey string, updateInterval time.Duration) *ConfigLoader {
	loader := &ConfigLoader{
		redisCli:       redisCli,
		pubsubKey:      pubsubKey,
		updateInterval: updateInterval,
		updateTasks:    map[string]func(){},
	}

	Go(func() {
		for {
			Safe(func() {
				pubsub := loader.redisCli.Subscribe(pubsubKey)
				defer pubsub.Close()
				for msg := range pubsub.Channel() {
					configFieled := msg.Payload
					if update, ok := loader.updateTasks[configFieled]; ok {
						update()
					}
				}
			})
		}
	})

	golog.Infof("NewConfigLoader, pubsubKey: %v, updateInterval: %v", pubsubKey, updateInterval)

	return loader
}

// 自动更新配置示例
// {
// 	// hset key
// 	ConfigUpdateKey := "CONFIG_UPDATE_KEY"
// 	// hset fieled
// 	ConfigUpdateFieled := "UPDATE_FIELED"
// 	// pubsub key
// 	ConfigUpdatePubsubKey := "CONFIG_UPDATE_PUBSUB_KEY"

// 	// 自动更新配置时间间隔
// 	autoUpdateInterval := time.Second * 15

// 	// 注意: 发布订阅会独占一个redis连接
// 	loader := util.NewConfigLoader(yabo.RedisCommon, ConfigUpdatePubsubKey, autoUpdateInterval)

// 	loader.Add(ConfigUpdateKey, ConfigUpdateFieled, onConfigUpdate)
// }