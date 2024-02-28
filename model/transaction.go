package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Id                uint64          `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	BlockNumber       uint64          `json:"blockNumber" gorm:"column:block_number;type:bigint(20);comment:块高度"`
	Hash              string          `json:"hash" gorm:"column:hash;type:varchar(200);comment:hash"`
	From              string          `json:"from" gorm:"column:from;type:varchar(80);comment:from"`
	To                string          `json:"to" gorm:"column:to;type:varchar(80);comment:to"`
	EffectiveGasPrice decimal.Decimal `json:"effectiveGasPrice" gorm:"column:effective_gas_price;type:decimal(30,0);comment:gas单价格"`
	GasUsed           uint64          `json:"gasUsed" gorm:"column:gas_used;type:bigint(20);comment:gas used"`
	TransactionIndex  uint            `json:"transactionIndex" gorm:"column:transaction_index;type:int(11);comment:交易索引"`
	ReceiptStatus     uint64          `json:"receiptStatus" gorm:"column:receipt_status;type:bigint(20);comment:交易回执状态"`
	Nonce             uint64          `json:"nonce" gorm:"column:nonce;type:bigint(20);comment:nonce"`
	Type              string          `json:"type" gorm:"column:type;type:char(1);comment:账变类型"`
	Value             decimal.Decimal `json:"value" gorm:"column:value;type:decimal(30,0);comment:账变数量"`
	Protected         string          `json:"protected" gorm:"column:protected;type:char(1);comment:是否保护(1-是 2-否)"`
	CreatedAt         *time.Time      `json:"createdAt" gorm:"column:created_at;type:datetime;comment:创建时间"`
}

func (Transaction) TableName() string {
	return "app_transaction"
}
