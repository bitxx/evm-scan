package config

import (
	"fmt"
	loadconfig "github.com/bitxx/load-config"
	"github.com/bitxx/load-config/source"
	"log"
)

type Config struct {
	Application *Application          `yaml:"application"`
	Logger      *Logger               `yaml:"logger"`
	Database    *Database             `yaml:"database"`
	Databases   *map[string]*Database `yaml:"databases"`
	Chain       *Chain                `yaml:"node"`
	callbacks   []func()
}

func (e *Config) Init() {
	e.Logger.Setup()
	e.multiDatabase()
	e.runCallback()
	log.Println("!!! client config init")
}

// 多db改造
func (e *Config) multiDatabase() {
	if len(*e.Databases) == 0 {
		*e.Databases = map[string]*Database{
			"*": e.Database,
		}

	}
}

func (e *Config) runCallback() {
	for i := range e.callbacks {
		e.callbacks[i]()
	}
}

func (e *Config) OnChange() {
	e.Init()
	log.Println("!!! client config change and reload")
}

// Setup 载入配置文件
func Setup(s source.Source,
	fs ...func()) {
	_cfg := &Config{
		Application: ApplicationConfig,
		Chain:       ChainConfig,
		Database:    DatabaseConfig,
		Databases:   &DatabasesConfig,
		Logger:      LoggerConfig,
		callbacks:   fs,
	}
	var err error
	loadconfig.DefaultConfig, err = loadconfig.NewConfig(
		loadconfig.WithSource(s),
		loadconfig.WithEntity(_cfg),
	)
	if err != nil {
		log.Println(fmt.Sprintf("New client config object fail: %s, use default param to start", err.Error()))
		return
	}
	_cfg.Init()
}
