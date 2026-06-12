package repository

import (
	"reporting-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(report *models.Report) error
	FindByID(id uuid.UUID) (*models.Report, error)
	FindAll(page, limit int) ([]models.Report, int64, error)
	FindByType(reportType string, page, limit int) ([]models.Report, int64, error)
	Delete(id uuid.UUID) error
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(report *models.Report) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) FindByID(id uuid.UUID) (*models.Report, error) {
	var report models.Report
	if err := r.db.First(&report, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) FindAll(page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.Report{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *reportRepository) FindByType(reportType string, page, limit int) ([]models.Report, int64, error) {
	var reports []models.Report
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&models.Report{}).Where("type = ?", reportType)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *reportRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Report{}, "id = ?", id).Error
}
