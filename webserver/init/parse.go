package init

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func parseLocalConfig(load loadFunc) {
	var configFile = ServerName + "-" + *runMode

	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("parse local config failed, err: " + err.Error())
	}

	var cfg = new(config)
	if err = viper.Unmarshal(cfg); err != nil {
		log.Println("unmarshal local config failed, err: " + err.Error())
	}

	if err = load(cfg); err != nil {
		log.Println("load local config failed, err: " + err.Error())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Printf("local config changed, event name[%s] value: %s\n", in.Name, in.String())

		var cfg = new(config)
		if err = viper.Unmarshal(cfg); err != nil {
			log.Println("unmarshal local config failed, err: " + err.Error())
		} else {
			if err = load(cfg); err != nil {
				log.Println("load local config failed, err: " + err.Error())
			}
		}
	})
}

func parseRemoteConfig(load loadFunc) {
	var err error
	if len(*secretkeyring) > 0 {
		err = viper.AddRemoteProvider(*provider, *endpoint, *wpath)
	} else {
		err = viper.AddSecureRemoteProvider(*provider, *endpoint, *wpath, *secretkeyring)
	}
	if err != nil {
		log.Println("parse remote config failed, err: " + err.Error())
	}

	viper.SetConfigType(*configType)
	if err = viper.ReadRemoteConfig(); err != nil {
		log.Println("read remote config failed, err: " + err.Error())
	}

	var cfg = new(config)
	if err = viper.Unmarshal(cfg); err != nil {
		log.Println("unmarshal remote config failed, err: " + err.Error())
	}

	if err = load(cfg); err != nil {
		log.Println("load remote config failed, err: " + err.Error())
	}

	go func() {
		for {
			time.Sleep(remoteConfigRefreshTime)

			var err error
			if err = viper.WatchRemoteConfig(); err != nil {
				log.Println("watch remote config failed, err: " + err.Error())
				continue
			}

			var cfg = new(config)
			if err = viper.Unmarshal(cfg); err != nil {
				log.Println("unmarshal remote config failed, err: " + err.Error())
				continue
			}

			if err = load(cfg); err != nil {
				log.Println("load remote config failed, err: " + err.Error())
			}
		}
	}()
}
