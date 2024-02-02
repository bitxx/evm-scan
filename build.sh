#!/usr/bin/env bash

go mod tidy

# macos arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o scan main.go
tar czvf "evm-inscriptions_-macos-arm64".tar.gz scan settings.yml app.sql README.md
rm -f scan

sleep 3

# macos amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o scan main.go
tar czvf "evm-inscriptions-macos-amd64".tar.gz scan settings.yml README.md
rm -f scan

sleep 3

# 交叉编译windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o scan.exe main.go
tar czvf "evm-inscriptions-windows".tar.gz scan.exe settings.yml README.md
rm -f scan.exe

sleep 3

# 交叉编译linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o scan main.go
tar czvf "evm-inscriptions-linux-amd64".tar.gz scan settings.yml README.md
rm -f scan
