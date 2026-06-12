package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reporting-service/internal/models"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type ExportService interface {
	ExportDashboardToCSV(metrics *models.DashboardMetrics) ([]byte, error)
	ExportDashboardToPDF(metrics *models.DashboardMetrics) ([]byte, error)
	ExportRevenueToCSV(report []models.RevenueReport) ([]byte, error)
	ExportRevenueToPDF(report []models.RevenueReport) ([]byte, error)
	ExportTopProductsToCSV(products []models.TopProduct) ([]byte, error)
	ExportTopProductsToPDF(products []models.TopProduct) ([]byte, error)
	ExportCustomersToCSV(customers []models.CustomerReport) ([]byte, error)
	ExportCustomersToPDF(customers []models.CustomerReport) ([]byte, error)
}

type exportService struct{}

func NewExportService() ExportService {
	return &exportService{}
}

// CSV Exports

func (s *exportService) ExportDashboardToCSV(metrics *models.DashboardMetrics) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Metric", "Value"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	records := [][]string{
		{"Total Revenue", fmt.Sprintf("$%.2f", metrics.TotalRevenue)},
		{"Total Orders", strconv.FormatInt(metrics.TotalOrders, 10)},
		{"Total Customers", strconv.FormatInt(metrics.TotalCustomers, 10)},
		{"Average Order Value", fmt.Sprintf("$%.2f", metrics.AverageOrderValue)},
		{"Today Revenue", fmt.Sprintf("$%.2f", metrics.TodayRevenue)},
		{"Today Orders", strconv.FormatInt(metrics.TodayOrders, 10)},
		{"Month Revenue", fmt.Sprintf("$%.2f", metrics.MonthRevenue)},
		{"Month Orders", strconv.FormatInt(metrics.MonthOrders, 10)},
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *exportService) ExportRevenueToCSV(report []models.RevenueReport) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Date", "Revenue", "Orders"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, r := range report {
		record := []string{
			r.Date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", r.Revenue),
			strconv.FormatInt(r.Orders, 10),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *exportService) ExportTopProductsToCSV(products []models.TopProduct) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Product ID", "Product Name", "Total Sold", "Total Revenue"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, p := range products {
		record := []string{
			p.ProductID.String(),
			p.ProductName,
			strconv.FormatInt(p.TotalSold, 10),
			fmt.Sprintf("%.2f", p.TotalRevenue),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

func (s *exportService) ExportCustomersToCSV(customers []models.CustomerReport) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"Customer ID", "Email", "Total Orders", "Total Spent", "Last Order Date"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, c := range customers {
		record := []string{
			c.CustomerID.String(),
			c.CustomerEmail,
			strconv.FormatInt(c.TotalOrders, 10),
			fmt.Sprintf("%.2f", c.TotalSpent),
			c.LastOrderDate.Format("2006-01-02"),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), nil
}

// PDF Exports

func (s *exportService) ExportDashboardToPDF(metrics *models.DashboardMetrics) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Dashboard Metrics Report")
	pdf.Ln(15)

	// Generated date
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(15)

	// Metrics table
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(100, 10, "Metric")
	pdf.Cell(0, 10, "Value")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	metrics_data := [][]string{
		{"Total Revenue", fmt.Sprintf("$%.2f", metrics.TotalRevenue)},
		{"Total Orders", strconv.FormatInt(metrics.TotalOrders, 10)},
		{"Total Customers", strconv.FormatInt(metrics.TotalCustomers, 10)},
		{"Average Order Value", fmt.Sprintf("$%.2f", metrics.AverageOrderValue)},
		{"Today Revenue", fmt.Sprintf("$%.2f", metrics.TodayRevenue)},
		{"Today Orders", strconv.FormatInt(metrics.TodayOrders, 10)},
		{"Month Revenue", fmt.Sprintf("$%.2f", metrics.MonthRevenue)},
		{"Month Orders", strconv.FormatInt(metrics.MonthOrders, 10)},
	}

	for _, row := range metrics_data {
		pdf.Cell(100, 8, row[0])
		pdf.Cell(0, 8, row[1])
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *exportService) ExportRevenueToPDF(report []models.RevenueReport) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Revenue Report")
	pdf.Ln(15)

	// Generated date
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(15)

	// Table header
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(60, 10, "Date")
	pdf.Cell(70, 10, "Revenue")
	pdf.Cell(0, 10, "Orders")
	pdf.Ln(10)

	// Table data
	pdf.SetFont("Arial", "", 11)
	for _, r := range report {
		pdf.Cell(60, 8, r.Date.Format("2006-01-02"))
		pdf.Cell(70, 8, fmt.Sprintf("$%.2f", r.Revenue))
		pdf.Cell(0, 8, strconv.FormatInt(r.Orders, 10))
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *exportService) ExportTopProductsToPDF(products []models.TopProduct) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape for wider table
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Top Products Report")
	pdf.Ln(15)

	// Generated date
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(15)

	// Table header
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(70, 10, "Product Name")
	pdf.Cell(40, 10, "Total Sold")
	pdf.Cell(0, 10, "Total Revenue")
	pdf.Ln(10)

	// Table data
	pdf.SetFont("Arial", "", 10)
	for _, p := range products {
		pdf.Cell(70, 8, p.ProductName)
		pdf.Cell(40, 8, strconv.FormatInt(p.TotalSold, 10))
		pdf.Cell(0, 8, fmt.Sprintf("$%.2f", p.TotalRevenue))
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *exportService) ExportCustomersToPDF(customers []models.CustomerReport) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape for wider table
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, "Customer Report")
	pdf.Ln(15)

	// Generated date
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(15)

	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(50, 10, "Customer ID")
	pdf.Cell(35, 10, "Total Orders")
	pdf.Cell(45, 10, "Total Spent")
	pdf.Cell(0, 10, "Last Order")
	pdf.Ln(10)

	// Table data
	pdf.SetFont("Arial", "", 9)
	for _, c := range customers {
		pdf.Cell(50, 8, c.CustomerID.String()[:18]+"...")
		pdf.Cell(35, 8, strconv.FormatInt(c.TotalOrders, 10))
		pdf.Cell(45, 8, fmt.Sprintf("$%.2f", c.TotalSpent))
		pdf.Cell(0, 8, c.LastOrderDate.Format("2006-01-02"))
		pdf.Ln(8)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
