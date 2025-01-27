package models

import "time"

type AggTradeTick struct {
	Symbol    string    `gorm:"primaryKey"                  json:"symbol"`
	Timestamp time.Time `gorm:"primaryKey;type:timestamptz" json:"timestamp"`
	Open      float64   `gorm:"not null"                    json:"open"`
	High      float64   `gorm:"not null"                    json:"high"`
	Low       float64   `gorm:"not null"                    json:"low"`
	Close     float64   `gorm:"not null"                    json:"close"`
	Volume    float64   `gorm:"not null"                    json:"volume"`
}

func (AggTradeTick) TableName() string {
	return "agg_trade_ticks"
}
