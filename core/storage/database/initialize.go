package database

import (
	"evm-scan/core/config"
	"evm-scan/core/runtime"
	"evm-scan/core/storage/database/cfg"
	"evm-scan/core/utils/log"
	"evm-scan/core/utils/textutils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

// Setup 配置数据库
func Setup() {
	for k := range config.DatabasesConfig {
		setupSimpleDatabase(k, config.DatabasesConfig[k])
	}
}

func setupSimpleDatabase(host string, c *config.Database) {
	registers := make([]cfg.ResolverConfigure, len(c.Registers))
	for i := range c.Registers {
		registers[i] = cfg.NewResolverConfigure(
			c.Registers[i].Sources,
			c.Registers[i].Replicas,
			c.Registers[i].Policy,
			c.Registers[i].Tables)
	}
	resolverConfig := cfg.NewConfigure(c.Source, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxIdleTime, c.ConnMaxLifeTime, registers)
	db, err := resolverConfig.Init(&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: New(
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					log.LevelForGorm()),
			},
		),
	}, opens[c.Driver])

	if err != nil {
		log.Fatal(textutils.Red(c.Driver+" connect error :"), err)
	} else {
		log.Info(textutils.Green(c.Driver + " connect success !"))
	}
	runtime.RuntimeConfig.SetDb(host, db)
}
