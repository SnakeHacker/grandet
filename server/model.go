package server

import "github.com/jinzhu/gorm"

// StockMeta ...
type StockMeta struct {
	gorm.Model
	TsCode   string `gorm:"column:ts_code; not null; unique"`
	Symbol   string `gorm:"column:symbol; not null; unique"`
	Name     string `gorm:"column:name; not null; unique"`
	Area     string `gorm:"column:area; not null"`
	Industry string `gorm:"column:industry; not null"`
	Market   string `gorm:"column:market; not null"`
	ListDate string `gorm:"column:list_date; not null"`
}

func (StockMeta) TableName() string {
	return "stocks"
}

// StockDaily ...
type StockDaily struct {
	gorm.Model
	TsCode    string  `gorm:"column:ts_code; not null"`
	TradeDate string  `gorm:"column:opentrade_date; not null"`
	Open      float32 `gorm:"column:open; not null"`
	High      float32 `gorm:"column:high; not null"`
	Low       float32 `gorm:"column:low; not null"`
	Close     float32 `gorm:"column:close; not null"`
	PreClose  float32 `gorm:"column:pre_close; not null"`
	Change    float32 `gorm:"column:change; not null"`
	PctChg    float32 `gorm:"column:pct_chg; not null"`
	Vol       int     `gorm:"column:vol; not null"`
	Amount    float32 `gorm:"column:amount; not null"`
}

func (StockDaily) TableName() string {
	return "daily"
}
