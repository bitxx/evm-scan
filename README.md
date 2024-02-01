# evm-scan
目前主要自用。可用于扫描evm链完整交易数据，里面也会陆续集成特定链的合约等数据，方便个人统计.  

## 当前功能
1. 支持evm链按块扫描交易，并录入mysql。线程池异步，针对网络情况，自己可调整配置参数，尽可能快的扫描交易。当前在`zkfair`网络测试通过  

## 编译步骤
1. 需要有golang环境，`1.21.5`及以上的版本  
2. 编译如下：
```shell
# 项目根目录添加依赖包
go mod tidy

# 直接编译当前平台
go build -o scan main.go

# 交叉编译windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o scan.exe main.go

# 交叉编译linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o scan main.go
```

## 使用步骤
1. 导入项目根目录的`app.sql`到mysql数据库，我用的`mariadb v10.7.3`  
2. 执行上一步骤编译好的程序（说明：配置文件中已经对各个配置做了详细描述）：
```shell
# 若不带任何参数，则程序读取当前目录的settings.yml配置文件
./scan

# 若指定配置文件，则按照指定文件读取
./scan settings.dev.yml
```
