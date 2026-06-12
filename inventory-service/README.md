# Inventory Service

**Port:** 8087  
**Language:** Golang  
**Framework:** Gin  
**Database:** PostgreSQL  
**ORM:** GORM

## Responsibilities

- **Stock Management** - Track product inventory
- **Stock Updates** - Update stock levels
- **Low Stock Alerts** - Notify when stock is low
- **Inventory Adjustments** - Manual stock adjustments
- **Stock History** - Track stock changes
- **Warehouse Management** - Multi-warehouse support

## Database Schema

### inventory
- id (UUID, PK)
- product_id (unique)
- sku
- quantity_available
- quantity_reserved
- quantity_sold
- warehouse_id (FK)
- low_stock_threshold
- created_at
- updated_at

### warehouses
- id (UUID, PK)
- name
- code (unique)
- address (JSON)
- is_active
- created_at
- updated_at

### inventory_transactions
- id (UUID, PK)
- inventory_id (FK)
- transaction_type (purchase, sale, adjustment, return)
- quantity
- previous_quantity
- new_quantity
- reference_id (order_id, etc.)
- notes
- created_at

## API Endpoints

- `GET /api/v1/inventory/:productId` - Get inventory for product
- `PUT /api/v1/inventory/:productId` - Update inventory (admin)
- `POST /api/v1/inventory/reserve` - Reserve stock (internal)
- `POST /api/v1/inventory/release` - Release reserved stock (internal)
- `GET /api/v1/inventory/low-stock` - Get low stock items (admin)
- `GET /api/v1/inventory/history/:productId` - Get stock history
- `POST /api/v1/warehouses` - Create warehouse (admin)
- `GET /api/v1/warehouses` - List warehouses

## Stock Management Flow

### Reserve Stock (on order creation)
```
quantity_available -= quantity
quantity_reserved += quantity
```

### Confirm Sale (on payment success)
```
quantity_reserved -= quantity
quantity_sold += quantity
```

### Release Stock (on order cancellation)
```
quantity_reserved -= quantity
quantity_available += quantity
```

## Low Stock Alerts

When `quantity_available < low_stock_threshold`:
- Publish `inventory.lowstock` event
- Notification service sends email to admin

## Events Published

- `inventory.updated`
- `inventory.lowstock`
- `inventory.outofstock`

## Events Consumed

- `order.created` - Reserve stock
- `order.confirmed` - Confirm sale
- `order.cancelled` - Release stock

## Directory Structure

```
inventory-service/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   └── events/
├── migrations/
├── pkg/
├── Dockerfile
├── go.mod
└── go.sum
```
