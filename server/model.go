package server

import (
	"github.com/jinzhu/gorm"
)

// StockMeta ...
type StockMeta struct {
	gorm.Model
	TsCode   string `gorm:"column:ts_code; not null; unique"`
	Symbol   string `gorm:"column:symbol; not null; unique"`
	Name     string `gorm:"column:name; not null; unique"`
	Area     string `gorm:"column:area; not null; index"`
	Industry string `gorm:"column:industry; not null; index"`
	Market   string `gorm:"column:market; not null; index"`
	ListDate string `gorm:"column:list_date; not null"`
}

func (StockMeta) TableName() string {
	return "stocks"
}

// StockDaily ...
type StockDaily struct {
	gorm.Model
	TsCode    string  `gorm:"column:ts_code; not null; index"`
	TradeDate string  `gorm:"column:trade_date; not null"`
	Open      float64 `gorm:"column:open; not null"`
	High      float64 `gorm:"column:high; not null"`
	Low       float64 `gorm:"column:low; not null"`
	Close     float64 `gorm:"column:close; not null"`
	PreClose  float64 `gorm:"column:pre_close; not null"`
	Change    float64 `gorm:"column:change; not null"`
	PctChg    float64 `gorm:"column:pct_chg; not null"`
	Vol       float64 `gorm:"column:vol; not null"`
	Amount    float64 `gorm:"column:amount; not null"`
}

func (StockDaily) TableName() string {
	return "daily"
}

// SortedDailys ...
type SortedDailys []StockDaily

func (p SortedDailys) Len() int { return len(p) }

func (p SortedDailys) Less(i, j int) bool {
	return p[i].Amount > p[j].Amount
}

func (p SortedDailys) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// ConceptDetail ...
type ConceptDetail struct {
	gorm.Model
	ConceptName string `gorm:"column:concept_name; not null"`
	TsCode      string `gorm:"column:ts_code; not null"`
}

func (ConceptDetail) TableName() string {
	return "concept_details"
}
