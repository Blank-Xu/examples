package init

import (
	"errors"
)

type loadFunc func(*config) error

func loadConfig(cfg *config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	cfg.Fix.Init()
	cfg.Log.Init()
	cfg.Server.Init()
	for _, v := range cfg.Jwt {
		v.Init()
	}

	return nil
}

// func Init() {
// 	flag.Parse()
//
// 	fmt.Printf("server version: %s, start args: %v\n", Version, flag.Args())
// 	fmt.Printf("load config file: %s\n", *configFile)
//
// 	file, err := os.OpenFile(*configFile, os.O_RDONLY, 0666)
// 	if err != nil {
// 		panic(fmt.Sprintf("open config file failed, err: %v", err))
// 	}
// 	defer file.Close()
//
// 	var cfg config
// 	if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
// 		panic(fmt.Sprintf("decode config file failed, err: %v", err))
// 	}
//
// 	// cfg.Fix.Init()
// 	// cfg.Log.Init()
// 	// cfg.Server.Init()
// 	// cfg.Jwt.Init()
//
// 	// log.SetOutput(io.MultiWriter(os.Stdout, zlog.Logger.Output()))
//
// 	Default = &cfg
//
// 	fmt.Println("load config file success")
// }
