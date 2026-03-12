package reports

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Report struct {
	ID          uint      `gorm:"primaryKey"`
	PeriodStart time.Time `gorm:"column:period_start"`
	PeriodEnd   time.Time `gorm:"column:period_end"`
	GeneratedAt time.Time `gorm:"column:generated_at"`
	Content     string    `gorm:"column:content"`
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Report, error) {
	var reports []Report
	err := r.db.WithContext(ctx).Order("generated_at DESC").Find(&reports).Error
	return reports, err
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*Report, error) {
	var report Report
	err := r.db.WithContext(ctx).First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *Repository) Save(ctx context.Context, report *Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}
