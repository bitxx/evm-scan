package app

import (
	"errors"
	"evm-scan/core/config"
	"evm-scan/core/runtime"
	"evm-scan/core/utils/log"
	appModel "evm-scan/model"
	"evm-scan/model/constant"
	"github.com/bitxx/evm-utils"
	"github.com/bitxx/evm-utils/model"
	"github.com/bitxx/evm-utils/util/dateutil"
	"gorm.io/gorm"
	"sync"
	"time"
)

type App struct {
	db      *gorm.DB
	client  *evmutils.EthClient
	txCache []appModel.Transaction
}

func NewApp() *App {

	return &App{
		db:      runtime.RuntimeConfig.GetDbByKey("*"),
		client:  evmutils.NewEthClient(config.ChainConfig.Url, config.ChainConfig.Timeout),
		txCache: []appModel.Transaction{},
	}
}

func (a *App) ScanAllTransactions() {
	blockNumStart := uint64(0)
	for {
		//需要延迟，以防某些异常导致无限请求
		time.Sleep(constant.TimeSleep)

		//获取块开始
		if blockNumStart <= 0 {
			blockNumStart = config.ChainConfig.BlockNumberStart
			tx := appModel.Transaction{}
			err := a.db.Last(&tx).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Errorf("txScan => get latest local block error: %s", err.Error())
				continue
			}
			// 此时如果err不为空，则说明ErrRecordNotFound，也就是 blockNumStart = 1
			if err == nil {
				blockNumStart = tx.BlockNumber + 1
			}
		}

		//获取块结束
		blockNumEnd, err := a.client.LatestBlockNumber()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("txScan => get latest chain block error: %s", err.Error())
			continue
		}

		//不断追最新块
		if blockNumStart < blockNumEnd {
			log.Infof("txScan => new scanning, block from %d to %d", blockNumStart, blockNumEnd)
			a.ScanTransactionsByNumber(blockNumStart, blockNumEnd)
			blockNumStart = blockNumEnd + 1
		}
	}

}

// ScanTransactionsByNumber
//
//	@Description: 扫描并添加交易记录
//	@receiver a
//	@param blockNumStart
//	@param blockNumEnd
func (a *App) ScanTransactionsByNumber(blockNumStart, blockNumEnd uint64) {
	txsCh := make(chan chan []appModel.Transaction, config.ChainConfig.BlockThreadSize)
	wg := sync.WaitGroup{}

	go a.saveTransactions(txsCh, blockNumStart, blockNumEnd)

	for i := blockNumStart; i <= blockNumEnd; i++ {
		time.Sleep(time.Duration(config.ChainConfig.BlockDelay) * time.Millisecond)

		c := make(chan []appModel.Transaction)
		txsCh <- c
		wg.Add(1)

		go func(i uint64) {
			var err error
			var transactions []model.Transaction
			for {
				transactions, err = a.client.TxReceiptByBlockNumber(i)
				if err != nil {
					log.Errorf("txScan => block: %d,transactions error: %s", i, err.Error())
					time.Sleep(constant.TimeSleep)
					continue
				}
				break
			}
			var ts []appModel.Transaction
			for _, tx := range transactions {
				protected := constant.TxProtectdFalse
				if tx.Protected {
					protected = constant.TxProtectdTrue
				}
				createTime, _ := dateutil.ParseTimestampToTime(int64(tx.Time), "UTC")
				t := appModel.Transaction{
					BlockNumber:       i,
					Hash:              tx.Hash,
					From:              tx.From,
					To:                tx.To,
					EffectiveGasPrice: tx.EffectiveGasPrice,
					GasUsed:           tx.GasUsed,
					TransactionIndex:  tx.TransactionIndex,
					ReceiptStatus:     tx.ReceiptStatus,
					Nonce:             tx.Nonce,
					Type:              tx.Type,
					Value:             tx.Value,

					Protected: protected,
					CreatedAt: &createTime,
				}
				ts = append(ts, t)
			}
			log.Infof("txScan => scanning ... block number: %d, tx count: %d", i, len(transactions))
			c <- ts
			wg.Done()
		}(i)
	}
	wg.Wait()

}

// saveTransactions
//
//	@Description: 保存交易记录
//	@receiver a
//	@param txsCh
//	@param blockNumStart
//	@param blockNumEnd
func (a *App) saveTransactions(txsCh chan chan []appModel.Transaction, blockNumStart, blockNumEnd uint64) {
	for {
		c := <-txsCh
		txs := <-c
		//批量插入
		a.txCache = append(a.txCache, txs...)
		if len(a.txCache) >= config.ChainConfig.TxCacheSize || uint64(len(a.txCache)) >= (blockNumEnd-blockNumStart)+1 {
			for {
				err := a.db.Create(a.txCache).Error
				if err != nil {
					log.Errorf("txScan => batch insert error, begin block: %d, end block: %d, error info: %s", a.txCache[0].BlockNumber, a.txCache[len(a.txCache)-1].BlockNumber, err)
					time.Sleep(constant.TimeSleep)
					continue
				}
				log.Infof("txScan => batch insert success, begin block: %d, end block: %d", a.txCache[0].BlockNumber, a.txCache[len(a.txCache)-1].BlockNumber)
				break
			}
			a.txCache = []appModel.Transaction{}
		}
	}

}
