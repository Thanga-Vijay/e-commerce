package handlers

import (
	"net/http"
	"reporting-service/internal/export"
	"reporting-service/internal/service"
	"reporting-service/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MetricsHandler struct {
	metricsService service.MetricsService
	reportService  service.ReportService
	exportService  export.ExportService
}

func NewMetricsHandler(metricsService service.MetricsService, reportService service.ReportService, exportService export.ExportService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
		reportService:  reportService,
		exportService:  exportService,
	}
}

// GetDashboard retrieves dashboard metrics
func (h *MetricsHandler) GetDashboard(c *gin.Context) {
	token, _ := c.Get("token")
	
	metrics, err := h.metricsService.GetDashboardMetrics(c.Request.Context(), token.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard metrics retrieved successfully", metrics)
}

// ExportDashboard exports dashboard metrics to CSV or PDF
func (h *MetricsHandler) ExportDashboard(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	token, _ := c.Get("token")
	
	metrics, err := h.metricsService.GetDashboardMetrics(c.Request.Context(), token.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var data []byte
	var contentType string
	var filename string

	switch format {
	case "csv":
		data, err = h.exportService.ExportDashboardToCSV(metrics)
		contentType = "text/csv"
		filename = "dashboard_metrics.csv"
	case "pdf":
		data, err = h.exportService.ExportDashboardToPDF(metrics)
		contentType = "application/pdf"
		filename = "dashboard_metrics.pdf"
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid format. Use 'csv' or 'pdf'")
		return
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// GetRevenueReport retrieves revenue report
func (h *MetricsHandler) GetRevenueReport(c *gin.Context) {
	token, _ := c.Get("token")
	
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	period := c.DefaultQuery("period", "daily")

	if startDateStr == "" || endDateStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "startDate and endDate are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid startDate format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid endDate format. Use YYYY-MM-DD")
		return
	}

	report, err := h.metricsService.GetRevenueReport(c.Request.Context(), token.(string), startDate, endDate, period)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Revenue report retrieved successfully", report)
}

// ExportRevenueReport exports revenue report to CSV or PDF
func (h *MetricsHandler) ExportRevenueReport(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	token, _ := c.Get("token")
	
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	period := c.DefaultQuery("period", "daily")

	if startDateStr == "" || endDateStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "startDate and endDate are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid startDate format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid endDate format")
		return
	}

	report, err := h.metricsService.GetRevenueReport(c.Request.Context(), token.(string), startDate, endDate, period)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var data []byte
	var contentType string
	var filename string

	switch format {
	case "csv":
		data, err = h.exportService.ExportRevenueToCSV(report)
		contentType = "text/csv"
		filename = "revenue_report.csv"
	case "pdf":
		data, err = h.exportService.ExportRevenueToPDF(report)
		contentType = "application/pdf"
		filename = "revenue_report.pdf"
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid format")
		return
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// GetTopProducts retrieves top products report
func (h *MetricsHandler) GetTopProducts(c *gin.Context) {
	token, _ := c.Get("token")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.metricsService.GetTopProducts(c.Request.Context(), token.(string), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Top products retrieved successfully", products)
}

// ExportTopProducts exports top products to CSV or PDF
func (h *MetricsHandler) ExportTopProducts(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	token, _ := c.Get("token")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.metricsService.GetTopProducts(c.Request.Context(), token.(string), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var data []byte
	var contentType string
	var filename string

	switch format {
	case "csv":
		data, err = h.exportService.ExportTopProductsToCSV(products)
		contentType = "text/csv"
		filename = "top_products.csv"
	case "pdf":
		data, err = h.exportService.ExportTopProductsToPDF(products)
		contentType = "application/pdf"
		filename = "top_products.pdf"
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid format")
		return
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// GetCustomerReport retrieves customer analytics
func (h *MetricsHandler) GetCustomerReport(c *gin.Context) {
	token, _ := c.Get("token")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	customers, err := h.metricsService.GetCustomerReports(c.Request.Context(), token.(string), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Customer report retrieved successfully", customers)
}

// ExportCustomerReport exports customer report to CSV or PDF
func (h *MetricsHandler) ExportCustomerReport(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	token, _ := c.Get("token")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	customers, err := h.metricsService.GetCustomerReports(c.Request.Context(), token.(string), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var data []byte
	var contentType string
	var filename string

	switch format {
	case "csv":
		data, err = h.exportService.ExportCustomersToCSV(customers)
		contentType = "text/csv"
		filename = "customer_report.csv"
	case "pdf":
		data, err = h.exportService.ExportCustomersToPDF(customers)
		contentType = "application/pdf"
		filename = "customer_report.pdf"
	default:
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid format")
		return
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}

// InvalidateCache invalidates all report caches
func (h *MetricsHandler) InvalidateCache(c *gin.Context) {
	if err := h.metricsService.InvalidateCache(c.Request.Context()); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cache invalidated successfully", nil)
}

// SaveReport saves a report to the database
func (h *MetricsHandler) SaveReport(c *gin.Context) {
	var req struct {
		Name      string    `json:"name" binding:"required"`
		Type      string    `json:"type" binding:"required"`
		Period    string    `json:"period"`
		StartDate *time.Time `json:"startDate"`
		EndDate   *time.Time `json:"endDate"`
		Data      interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, _ := c.Get("userID")
	
	report, err := h.reportService.SaveReport(req.Name, req.Type, req.Period, req.StartDate, req.EndDate, req.Data, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Report saved successfully", report)
}

// GetSavedReports retrieves all saved reports
func (h *MetricsHandler) GetSavedReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reports, total, err := h.reportService.GetAllReports(page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := gin.H{
		"reports": reports,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Reports retrieved successfully", response)
}

// GetReportByID retrieves a saved report by ID
func (h *MetricsHandler) GetReportByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid report ID")
		return
	}

	report, err := h.reportService.GetReportByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Report not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Report retrieved successfully", report)
}

// DeleteReport deletes a saved report
func (h *MetricsHandler) DeleteReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid report ID")
		return
	}

	if err := h.reportService.DeleteReport(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Report deleted successfully", nil)
}
