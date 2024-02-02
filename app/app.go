package app

import (
	"errors"
	"evm-scan/core/config"
	"evm-scan/core/runtime"
	"evm-scan/core/utils/log"
	appModel "evm-scan/model"
	"evm-scan/model/constant"
	"fmt"
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
	for {
		//获取块开始
		tx := appModel.Transaction{}
		err := a.db.Last(&tx).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(fmt.Sprintf("get latest local block error: %s", err.Error()))
			continue
		}
		blockNumStart := int64(1)
		if err == nil {
			blockNumStart = tx.BlockNumber + 1
		}

		//获取块结束
		blockNumEnd, err := a.client.LatestBlockNumber()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(fmt.Sprintf("get latest chain block error: %s", err.Error()))
			continue
		}

		//不断追最新块
		if uint64(blockNumStart) < blockNumEnd {
			a.ScanTransactionsByNumber(uint64(blockNumStart), blockNumEnd)
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
				transactions, err = a.client.TransactionsByBlockNumber(i)
				if err != nil {
					log.Error(fmt.Sprintf("block: %d,transactions error: %s", i, err.Error()))
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
				createTime := dateutil.ParseTimestampToTime(int64(tx.Time))
				t := appModel.Transaction{
					Hash:        tx.Hash,
					BlockNumber: int64(i),
					From:        tx.From,
					To:          tx.To,
					GasPrice:    tx.GasPrice,
					Cost:        tx.Cost,
					Type:        tx.Type,
					Value:       tx.Value,
					Protected:   protected,
					CreatedAt:   &createTime,
				}
				ts = append(ts, t)
			}
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
		for idx, tx := range txs {
			log.Infof("scanning ... block number: %d, tx number: %d", tx.BlockNumber, idx+1)
		}
		//批量插入
		a.txCache = append(a.txCache, txs...)
		if len(a.txCache) >= config.ChainConfig.TxCacheSize || blockNumEnd-blockNumStart < uint64(config.ChainConfig.TxCacheSize) {
			for {
				err := a.db.Create(a.txCache).Error
				if err != nil {
					log.Errorf("batch insert error,begin block: %d, end block: %d, begin tx: %s, end tx: %s, error info", a.txCache[0].BlockNumber, a.txCache[len(a.txCache)-1].BlockNumber, a.txCache[0].Hash, a.txCache[len(a.txCache)-1].Hash, err)
					time.Sleep(constant.TimeSleep)
					continue
				}
				break
			}
			a.txCache = []appModel.Transaction{}
		}
	}

}
