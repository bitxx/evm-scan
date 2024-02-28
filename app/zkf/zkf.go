package zkf

import (
	"context"
	"errors"
	zkfModel "evm-scan/app/zkf/model"
	zkfConstant "evm-scan/app/zkf/model/constant"
	"evm-scan/core/config"
	"evm-scan/core/runtime"
	"evm-scan/core/utils/log"
	appModel "evm-scan/model"
	"evm-scan/model/constant"
	"github.com/bitxx/evm-utils"
	"github.com/bitxx/evm-utils/util/dateutil"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type ZKF struct {
	db        *gorm.DB
	client    *evmutils.EthClient
	cacheScan map[string]uint64
}

func NewZKF() *ZKF {

	return &ZKF{
		db:        runtime.RuntimeConfig.GetDbByKey("*"),
		client:    evmutils.NewEthClient(config.ChainConfig.Url, config.ChainConfig.Timeout),
		cacheScan: map[string]uint64{},
	}
}

func (z *ZKF) StatGas() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	go z.statGasByTableName(zkfModel.HourTable)
	go z.statGasByTableName(zkfModel.DailyTable)
	go z.statGasByTableName(zkfModel.WeeklyTable)
	<-ctx.Done()
}

// statGasByTableName
//
//	@Description: 根据表统计gas情况
//	@receiver z
func (z *ZKF) statGasByTableName(tableName string) {
	if tableName != zkfModel.HourTable && tableName != zkfModel.DailyTable && tableName != zkfModel.WeeklyTable {
		log.Error("zkf => table name %s is error", tableName)
		return
	}
	//UTC启动时间
	dateStart, err := dateutil.ParseStrToTime("2024-01-14 00:00:00", "UTC", -1)
	if err != nil {
		log.Error("zkf => get default start date error: ", err)
		return
	}
	dateEnd := dateStart
	for {
		//需要延迟，以防某些异常导致无限请求
		time.Sleep(constant.TimeSleep)

		//减少误差，当前时间在此处标记
		now := time.Now()

		zkfStat := zkfModel.ZkfStatGas{}
		err := z.db.Table(tableName).Last(&zkfStat).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("zkf => [%s] get latest gas stat info error: %s", tableName, err.Error())
			continue
		}

		// 此时如果err不为空，则说明ErrRecordNotFound，日期从默认值开始算
		//若为空，则表明有记录
		if err == nil {
			//默认上一个阶段的已结束，则开启新的时间
			dateStart = zkfStat.DateEnd.Add(1 * time.Second)
			if zkfStat.CalcStatus == zkfConstant.CalcStatusRunning {
				dateStart = *zkfStat.DateStart
			}
		}

		// 读取一个阶段的统计结果
		var totalGasFee decimal.Decimal
		var totalGasPrice decimal.Decimal
		var minBlockNumber uint64
		var maxBlockNumber uint64
		var maxBlockNumberDate time.Time
		var count int64

		calcStatus := zkfConstant.CalcStatusRunning
		flag := "hour"
		if tableName == zkfModel.HourTable {
			flag = "hour"
			dateEnd = dateStart.Add(3599 * time.Second)

			//统计
			row := z.db.Model(&appModel.Transaction{}).Select("IFNULL(sum(gas_used * effective_gas_price),0),IFNULL(sum(effective_gas_price),0),IFNULL(min(block_number),0),IFNULL(max(block_number),0),IFNULL(max(created_at),now()),IFNULL(count(1),0)").Where("created_at>=? and created_at<=?", dateStart, dateEnd).Row()
			err = row.Scan(&totalGasFee, &totalGasPrice, &minBlockNumber, &maxBlockNumber, &maxBlockNumberDate, &count)
			if err != nil {
				log.Error("zkf => [%s] stat gas data err: %s", flag, err.Error())
				continue
			}
			if count <= 0 {
				log.Warnf("zkf => [%s] tx data is empty", flag)
				continue
			}

			//若没有新增数据，则继续下一轮
			if z.cacheScan[tableName] == maxBlockNumber {
				continue
			}
			z.cacheScan[tableName] = maxBlockNumber

			//判断这个阶段的数据是否结束
			tx := appModel.Transaction{}
			err = z.db.Last(&tx).Error
			if err != nil {
				log.Error("zkf => [%s] get latest tx err: %s", flag, err.Error())
				continue
			}
			if tx.CreatedAt.Compare(dateEnd) > 0 {
				calcStatus = zkfConstant.CalcStatusStop
			}
		} else {
			tbName := zkfModel.HourTable
			if tableName == zkfModel.DailyTable {
				flag = "daily"
				tbName = zkfModel.HourTable
				dateEnd = dateStart.AddDate(0, 0, 1).Add(-1 * time.Second)
			}
			if tableName == zkfModel.WeeklyTable {
				flag = "weekly"
				tbName = zkfModel.DailyTable
				dateEnd = dateStart.AddDate(0, 0, 7).Add(-1 * time.Second)
			}
			//统计
			row := z.db.Table(tbName).Select("IFNULL(sum(total_gas_fee),0),IFNULL(sum(total_gas_price),0),IFNULL(min(block_start),0),IFNULL(max(block_end),0),IFNULL(max(date_end),now()),IFNULL(sum(total_tx_count),0)").Where("date_start>=? and date_end<=?", dateStart, dateEnd).Row()
			err = row.Scan(&totalGasFee, &totalGasPrice, &minBlockNumber, &maxBlockNumber, &maxBlockNumberDate, &count)
			if err != nil {
				log.Error("zkf => [%s] stat gas data err: %s", flag, err.Error())
				continue
			}

			if count <= 0 {
				log.Warnf("zkf => [%s] %s data is empty", flag, tbName)
				continue
			}

			//若没有新增数据，则继续下一轮
			if z.cacheScan[tableName] == maxBlockNumber {
				continue
			}
			z.cacheScan[tableName] = maxBlockNumber

			//判断这个阶段的数据是否结束
			zkfStatGas := zkfModel.ZkfStatGas{}
			err = z.db.Table(tbName).Last(&zkfStatGas).Error
			if err != nil {
				log.Errorf("zkf => [%s] get latest tx err: %s", flag, err.Error())
				continue
			}
			if zkfStatGas.DateEnd.Compare(dateEnd) > 0 {
				calcStatus = zkfConstant.CalcStatusStop
			}
		}

		result := zkfModel.ZkfStatGas{}
		result.DateStart = &dateStart
		result.DateEnd = &maxBlockNumberDate
		//若结束，则用截止时间
		if calcStatus == zkfConstant.CalcStatusStop {
			result.DateEnd = &dateEnd
		}
		result.BlockStart = minBlockNumber
		result.BlockEnd = maxBlockNumber
		result.TotalGasFee = totalGasFee
		result.TotalGasPrice = totalGasPrice
		result.TotalTxCount = count
		result.CalcStatus = calcStatus

		if zkfStat.CalcStatus == zkfConstant.CalcStatusRunning && zkfStat.Id > 0 {
			result.Id = zkfStat.Id
			result.UpdatedAt = &now
			result.CreatedAt = zkfStat.CreatedAt
		} else {
			result.CreatedAt = &now
			result.UpdatedAt = &now
		}
		z.db.Table(tableName).Save(&result)

		log.Infof("zkf => calc over [%s] start block %d, end block %d", flag, minBlockNumber, maxBlockNumber)
	}
}
