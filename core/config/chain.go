package config

type Chain struct {
	Url              string
	Timeout          int64
	BlockDelay       int
	TxCacheSize      int
	BlockThreadSize  int
	BlockNumberStart uint64
}

var ChainConfig = new(Chain)
