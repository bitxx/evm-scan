package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Id          int64           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	BlockNumber int64           `json:"blockNumber" gorm:"column:block_number;type:int(11);comment:块高度"`
	Hash        string          `json:"hash" gorm:"column:hash;type:varchar(200);comment:hash"`
	From        string          `json:"from" gorm:"column:from;type:varchar(80);comment:from"`
	To          string          `json:"to" gorm:"column:to;type:varchar(80);comment:to"`
	GasPrice    decimal.Decimal `json:"gasPrice" gorm:"column:gas_price;type:decimal(18,6);comment:gas单价格"`
	Cost        decimal.Decimal `json:"cost" gorm:"column:cost;type:decimal(18,6);comment:gas费"`
	Type        string          `json:"type" gorm:"column:type;type:char(1);comment:账变类型"`
	Value       decimal.Decimal `json:"value" gorm:"column:value;type:decimal(18,6);comment:账变数量"`
	Protected   string          `json:"protected" gorm:"column:protected;type:char(1);comment:是否保护(1-是 2-否)"`
	CreatedAt   *time.Time      `json:"createdAt" gorm:"column:created_at;type:datetime;comment:创建时间"`
}

func (Transaction) TableName() string {
	return "app_transaction"
}
