package aggtrade

import (
	"context"
	"fmt"
	"github.com/majidmvulle/binance-trading-chart-service/persistor/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) SaveTick(ctx context.Context, tick models.AggTradeTick) error {
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "symbol"}, {Name: "timestamp"}},
		DoUpdates: clause.AssignmentColumns([]string{"open", "high", "low", "close", "volume"}),
	}).Create(&tick)

	if result.Error != nil {
		return fmt.Errorf("failed to save candlestick to database: %w", result.Error)
	}

	fmt.Printf("received Candlestick: Symbol=%s, Timestamp=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.2f\n",
		tick.Symbol, tick.Timestamp.Format(time.RFC3339), tick.Open, tick.High, tick.Low, tick.Close, tick.Volume)

	return nil
}
