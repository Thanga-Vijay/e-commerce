package service

import (
	"encoding/json"
	"fmt"
	"reporting-service/internal/models"
	"reporting-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ReportService interface {
	SaveReport(name, reportType, period string, startDate, endDate *time.Time, data interface{}, generatedBy uuid.UUID) (*models.Report, error)
	GetReportByID(id uuid.UUID) (*models.Report, error)
	GetAllReports(page, limit int) ([]models.Report, int64, error)
	GetReportsByType(reportType string, page, limit int) ([]models.Report, int64, error)
	DeleteReport(id uuid.UUID) error
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) SaveReport(name, reportType, period string, startDate, endDate *time.Time, data interface{}, generatedBy uuid.UUID) (*models.Report, error) {
	// Convert data to JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	report := &models.Report{
		Name:        name,
		Type:        reportType,
		Period:      period,
		StartDate:   startDate,
		EndDate:     endDate,
		Data:        string(jsonData),
		GeneratedBy: generatedBy,
	}

	if err := s.repo.Create(report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return report, nil
}

func (s *reportService) GetReportByID(id uuid.UUID) (*models.Report, error) {
	return s.repo.FindByID(id)
}

func (s *reportService) GetAllReports(page, limit int) ([]models.Report, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindAll(page, limit)
}

func (s *reportService) GetReportsByType(reportType string, page, limit int) ([]models.Report, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.repo.FindByType(reportType, page, limit)
}

func (s *reportService) DeleteReport(id uuid.UUID) error {
	return s.repo.Delete(id)
}
