package models

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	HourTable   = "app_zkf_stat_hours_gas"
	DailyTable  = "app_zkf_stat_daily_gas"
	WeeklyTable = "app_zkf_stat_weekly_gas"
)

var tableName = ""

type ZkfStatGas struct {
	Id           int64           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	BlockStart   int64           `json:"blockStart" gorm:"column:block_start;type:int(11);comment:块开始高度"`
	BlockEnd     int64           `json:"blockEnd" gorm:"column:block_end;type:int(11);comment:块截止高度"`
	DateStart    *time.Time      `json:"dateStart" gorm:"column:date_start;type:datetime;comment:开始时间"`
	DateEnd      *time.Time      `json:"dateEnd" gorm:"column:date_end;type:datetime;comment:截止时间"`
	TotalTxCount int64           `json:"totalTxCount" gorm:"column:total_tx_count;type:int(11);comment:交易总笔数"`
	TotalGas     decimal.Decimal `json:"totalGas" gorm:"column:total_gas;type:decimal(30,0);comment:交易总gas"`
	AvgGasFee    decimal.Decimal `json:"avgGasFee" gorm:"column:avg_gas_fee;type:decimal(30,0);comment:平均每笔gas"`
	AvgGasPrice  decimal.Decimal `json:"avgGasPrice" gorm:"column:avg_gas_price;type:decimal(30,0);comment:平均每笔gas price"`
	CalcStatus   string          `json:"calcStatus" gorm:"column:calc_status;type:char(1);comment:状态(1-运行中 2-停止)"`
	UpdatedAt    *time.Time      `json:"updatedAt" gorm:"column:updated_at;type:datetime;comment:更新时间"`
	CreatedAt    *time.Time      `json:"createdAt" gorm:"column:created_at;type:datetime;comment:创建时间"`
}
