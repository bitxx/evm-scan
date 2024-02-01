package main

import (
	"context"
	"evm-scan/app"
	"evm-scan/core/config"
	"evm-scan/core/storage/database"
	"evm-scan/core/utils/log"
	"evm-scan/model/constant"
	"github.com/bitxx/load-config/source/file"
	"os"
	"strings"
	"time"
)

func init() {
	configPath := "settings.yml"
	if len(os.Args) >= 2 {
		configPath = os.Args[1]
	}
	if !strings.HasSuffix(configPath, ".yml") {
		panic("config file error,please check it.")
	}
	config.Setup(
		file.NewSource(file.WithPath(configPath)),
		database.Setup,
	)
}

func main() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if err != nil {
			log.Errorf("执行异常，程序终止，失败原因：%s", err)
			return
		}
		cancel()
		log.Info("操作成功！")
		time.Sleep(constant.TimeSleep) //3秒种后，窗口关闭
	}()

	go func() {
		app.NewApp().ScanAllTransactions()
	}()

	<-ctx.Done()
}
