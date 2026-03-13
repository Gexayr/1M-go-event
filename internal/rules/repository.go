package rules

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]RiskRule, error)
	GetByID(id uint) (*RiskRule, error)
	Create(rule *RiskRule) error
	Update(rule *RiskRule) error
	Delete(id uint) error
	GetEnabledRules() ([]RiskRule, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]RiskRule, error) {
	var rules []RiskRule
	err := r.db.Find(&rules).Error
	return rules, err
}

func (r *repository) GetByID(id uint) (*RiskRule, error) {
	var rule RiskRule
	err := r.db.First(&rule, id).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *repository) Create(rule *RiskRule) error {
	return r.db.Create(rule).Error
}

func (r *repository) Update(rule *RiskRule) error {
	return r.db.Save(rule).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&RiskRule{}, id).Error
}

func (r *repository) GetEnabledRules() ([]RiskRule, error) {
	var rules []RiskRule
	err := r.db.Where("enabled = ?", true).Find(&rules).Error
	return rules, err
}
