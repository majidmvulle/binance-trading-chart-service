package aggtrade

import (
	"context"
	"fmt"

	"github.com/majidmvulle/binance-trading-chart-service/persistor/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	return nil
}
