package templates

import (
	"bytes"
	"fmt"
	"html/template"
)

// Template data structures
type WelcomeData struct {
	Name string
}

type EmailVerificationData struct {
	Name             string
	VerificationLink string
}

type PasswordResetData struct {
	Name      string
	ResetLink string
}

type OrderConfirmationData struct {
	Name        string
	OrderNumber string
	Total       float64
	Items       []OrderItem
}

type OrderItem struct {
	Name     string
	Quantity int
	Price    float64
}

type PaymentReceiptData struct {
	Name        string
	OrderNumber string
	Amount      float64
	PaymentDate string
}

type OrderShippedData struct {
	Name           string
	OrderNumber    string
	TrackingNumber string
	Carrier        string
}

type OrderDeliveredData struct {
	Name        string
	OrderNumber string
}

type LowStockAlertData struct {
	ProductName string
	SKU         string
	Quantity    int
	Threshold   int
}

// GetTemplate returns the HTML template for a given template name
func GetTemplate(templateName string) (*template.Template, error) {
	templates := map[string]string{
		"welcome":              welcomeTemplate,
		"email_verification":   emailVerificationTemplate,
		"password_reset":       passwordResetTemplate,
		"order_confirmation":   orderConfirmationTemplate,
		"payment_receipt":      paymentReceiptTemplate,
		"order_shipped":        orderShippedTemplate,
		"order_delivered":      orderDeliveredTemplate,
		"low_stock_alert":      lowStockAlertTemplate,
	}

	tmplStr, exists := templates[templateName]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	tmpl, err := template.New(templateName).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with the provided data
func RenderTemplate(templateName string, data interface{}) (string, error) {
	tmpl, err := GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// Template strings
const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to E-Commerce Platform</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #4CAF50;">Welcome to E-Commerce Platform!</h1>
        <p>Hi {{.Name}},</p>
        <p>Thank you for joining E-Commerce Platform. We're excited to have you on board!</p>
        <p>Start exploring our wide range of products and enjoy a seamless shopping experience.</p>
        <p>If you have any questions, feel free to reach out to our support team.</p>
        <p>Happy shopping!</p>
        <p style="margin-top: 30px;">Best regards,<br>The E-Commerce Team</p>
    </div>
</body>
</html>
`

const emailVerificationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #4CAF50;">Verify Your Email Address</h1>
        <p>Hi {{.Name}},</p>
        <p>Please verify your email address by clicking the button below:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.VerificationLink}}" style="background-color: #4CAF50; color: white; padding: 12px 30px; text-decoration: none; border-radius: 4px; display: inline-block;">Verify Email</a>
        </div>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #666;">{{.VerificationLink}}</p>
        <p>This link will expire in 24 hours.</p>
        <p style="margin-top: 30px;">Best regards,<br>The E-Commerce Team</p>
    </div>
</body>
</html>
`

const passwordResetTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #4CAF50;">Reset Your Password</h1>
        <p>Hi {{.Name}},</p>
        <p>We received a request to reset your password. Click the button below to set a new password:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{.ResetLink}}" style="background-color: #4CAF50; color: white; padding: 12px 30px; text-decoration: none; border-radius: 4px; display: inline-block;">Reset Password</a>
        </div>
        <p>Or copy and paste this link into your browser:</p>
        <p style="word-break: break-all; color: #666;">{{.ResetLink}}</p>
        <p>If you didn't request a password reset, you can safely ignore this email.</p>
        <p>This link will expire in 1 hour.</p>
        <p style="margin-top: 30px;">Best regards,<br>The E-Commerce Team</p>
    </div>
</body>
</html>
`

const orderConfirmationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Confirmation</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #4CAF50;">Order Confirmed!</h1>
        <p>Hi {{.Name}},</p>
        <p>Thank you for your order! Your order has been confirmed and is being processed.</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <p><strong>Order Number:</strong> {{.OrderNumber}}</p>
            <p><strong>Total:</strong> ${{printf "%.2f" .Total}}</p>
        </div>
        <h3>Order Items:</h3>
        <table style="width: 100%; border-collapse: collapse;">
            {{range .Items}}
            <tr style="border-bottom: 1px solid #ddd;">
                <td style="padding: 10px;">{{.Name}}</td>
                <td style="padding: 10px; text-align: center;">x{{.Quantity}}</td>
                <td style="padding: 10px; text-align: right;">${{printf "%.2f" .Price}}</td>
            </tr>
            {{end}}
        </table>
        <p style="margin-top: 20px;">We'll send you another email when your order ships.</p>
        <p style="margin-top: 30px;">Best regards,<br>The E-Commerce Team</p>
    </div>
</body>
</html>
`

const paymentReceiptTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Payment Receipt</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #4CAF50;">Payment Received</h1>
        <p>Hi {{.Name}},</p>
        <p>We've received your payment. Thank you!</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <p><strong>Order Number:</strong> {{.OrderNumber}}</p>
            <p><strong>Amount:</strong> ${{printf "%.2f" .Amount}}</p>
            <p><strong>Payment Date:</strong> {{.PaymentDate}}</p>
        </div>
        <p>Your order is now being prepared for shipment.</p>
        <p style="margin-top: 30px;">Best regards,<br>The E-Commerce Team</p>
    </div>
</body>
</html>
`

const orderShippedTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Shipped</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #4CAF50;">Your Order Has Shipped!</h1>
        <p>Hi {{.Name}},</p>
        <p>Great news! Your order has been shipped and is on its way to you.</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <p><strong>Order Number:</strong> {{.OrderNumber}}</p>
            <p><strong>Tracking Number:</strong> {{.TrackingNumber}}</p>
            <p><strong>Carrier:</strong> {{.Carrier}}</p>
        </div>
        <p>You can track your package using the tracking number provided above.</p>
        <p style="margin-top: 30px;">Best regards,<br>The E-Commerce Team</p>
    </div>
</body>
</html>
`

const orderDeliveredTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Order Delivered</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #4CAF50;">Order Delivered!</h1>
        <p>Hi {{.Name}},</p>
        <p>Your order has been successfully delivered!</p>
        <div style="background-color: #f5f5f5; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <p><strong>Order Number:</strong> {{.OrderNumber}}</p>
        </div>
        <p>We hope you're happy with your purchase. If you have any questions or concerns, please don't hesitate to contact us.</p>
        <p>Thank you for shopping with us!</p>
        <p style="margin-top: 30px;">Best regards,<br>The E-Commerce Team</p>
    </div>
</body>
</html>
`

const lowStockAlertTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Low Stock Alert</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #FF5722;">Low Stock Alert</h1>
        <p>This is an automated alert to notify you that a product is running low on stock.</p>
        <div style="background-color: #fff3e0; padding: 15px; border-left: 4px solid #FF5722; margin: 20px 0;">
            <p><strong>Product:</strong> {{.ProductName}}</p>
            <p><strong>SKU:</strong> {{.SKU}}</p>
            <p><strong>Current Quantity:</strong> {{.Quantity}}</p>
            <p><strong>Threshold:</strong> {{.Threshold}}</p>
        </div>
        <p>Please reorder this product to avoid stockouts.</p>
        <p style="margin-top: 30px;">Best regards,<br>Inventory Management System</p>
    </div>
</body>
</html>
`
