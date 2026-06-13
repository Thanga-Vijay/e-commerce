# Database Entity-Relationship Diagrams

Complete database schema design for all microservices in the E-Commerce Platform.

## Table of Contents
1. [Auth Service Database](#auth-service-database)
2. [Product Service Database](#product-service-database)
3. [Cart Service Database](#cart-service-database)
4. [Wishlist Service Database](#wishlist-service-database)
5. [Order Service Database](#order-service-database)
6. [Payment Service Database](#payment-service-database)
7. [Inventory Service Database](#inventory-service-database)
8. [Notification Service Database](#notification-service-database)
9. [Reporting Service Database](#reporting-service-database)
10. [Database Indexes](#database-indexes)
11. [Migration Strategy](#migration-strategy)

---

## Auth Service Database

**Database Name:** `auth_db`

### ER Diagram

```mermaid
erDiagram
    users ||--o{ refresh_tokens : has
    users ||--o{ password_resets : has
    
    users {
        uuid id PK
        varchar email UK
        varchar password_hash
        varchar first_name
        varchar last_name
        varchar role
        boolean is_verified
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    
    refresh_tokens {
        uuid id PK
        uuid user_id FK
        varchar token UK
        timestamp expires_at
        timestamp created_at
    }
    
    password_resets {
        uuid id PK
        uuid user_id FK
        varchar token UK
        timestamp expires_at
        boolean used
        timestamp created_at
    }
```

### Table Descriptions

#### users
Primary user authentication table.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique user identifier |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User email address |
| password_hash | VARCHAR(255) | NOT NULL | bcrypt hashed password |
| first_name | VARCHAR(100) | | User's first name |
| last_name | VARCHAR(100) | | User's last name |
| role | VARCHAR(20) | DEFAULT 'customer' | User role (customer, admin) |
| is_verified | BOOLEAN | DEFAULT false | Email verification status |
| created_at | TIMESTAMP | NOT NULL | Account creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |
| deleted_at | TIMESTAMP | | Soft delete timestamp |

#### refresh_tokens
Manages JWT refresh tokens for session management.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Token identifier |
| user_id | UUID | FOREIGN KEY (users.id), NOT NULL | Associated user |
| token | VARCHAR(500) | UNIQUE, NOT NULL | Refresh token string |
| expires_at | TIMESTAMP | NOT NULL | Token expiration |
| created_at | TIMESTAMP | NOT NULL | Token creation time |

#### password_resets
Tracks password reset requests.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Reset request identifier |
| user_id | UUID | FOREIGN KEY (users.id), NOT NULL | User requesting reset |
| token | VARCHAR(255) | UNIQUE, NOT NULL | Reset token |
| expires_at | TIMESTAMP | NOT NULL | Token expiration (1 hour) |
| used | BOOLEAN | DEFAULT false | Whether token was used |
| created_at | TIMESTAMP | NOT NULL | Request creation time |

---

## Product Service Database

**Database Name:** `product_db`

### ER Diagram

```mermaid
erDiagram
    categories ||--o{ categories : "parent-child"
    categories ||--o{ products : contains
    products ||--o{ reviews : has
    
    products {
        uuid id PK
        varchar name
        text description
        decimal price
        uuid category_id FK
        varchar brand
        varchar sku UK
        text[] images
        decimal average_rating
        int total_reviews
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
    
    categories {
        uuid id PK
        varchar name UK
        varchar slug UK
        uuid parent_id FK
        text description
        timestamp created_at
        timestamp updated_at
    }
    
    reviews {
        uuid id PK
        uuid product_id FK
        uuid user_id
        int rating
        varchar title
        text comment
        boolean verified_purchase
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
```

### Table Descriptions

#### products
Product catalog information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Product identifier |
| name | VARCHAR(255) | NOT NULL | Product name |
| description | TEXT | | Detailed product description |
| price | DECIMAL(10,2) | NOT NULL | Product price |
| category_id | UUID | FOREIGN KEY (categories.id) | Product category |
| brand | VARCHAR(100) | | Product brand |
| sku | VARCHAR(100) | UNIQUE, NOT NULL | Stock keeping unit |
| images | TEXT[] | | Array of image URLs |
| average_rating | DECIMAL(3,2) | DEFAULT 0 | Average customer rating (0-5) |
| total_reviews | INT | DEFAULT 0 | Total number of reviews |
| created_at | TIMESTAMP | NOT NULL | Product creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |
| deleted_at | TIMESTAMP | | Soft delete timestamp |

#### categories
Hierarchical product categories.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Category identifier |
| name | VARCHAR(100) | UNIQUE, NOT NULL | Category name |
| slug | VARCHAR(100) | UNIQUE, NOT NULL | URL-friendly slug |
| parent_id | UUID | FOREIGN KEY (categories.id), NULLABLE | Parent category |
| description | TEXT | | Category description |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

#### reviews
Customer product reviews and ratings.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Review identifier |
| product_id | UUID | FOREIGN KEY (products.id), NOT NULL | Reviewed product |
| user_id | UUID | NOT NULL | Reviewer user ID |
| rating | INT | CHECK (rating >= 1 AND rating <= 5) | Star rating (1-5) |
| title | VARCHAR(200) | | Review title |
| comment | TEXT | | Review text |
| verified_purchase | BOOLEAN | DEFAULT false | Verified purchase indicator |
| created_at | TIMESTAMP | NOT NULL | Review creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |
| deleted_at | TIMESTAMP | | Soft delete timestamp |

---

## Cart Service Database

**Database Name:** `cart_db`

### ER Diagram

```mermaid
erDiagram
    carts ||--o{ cart_items : contains
    
    carts {
        uuid id PK
        uuid user_id UK
        timestamp created_at
        timestamp updated_at
        timestamp expires_at
    }
    
    cart_items {
        uuid id PK
        uuid cart_id FK
        uuid product_id
        int quantity
        decimal price
        timestamp created_at
        timestamp updated_at
    }
```

### Table Descriptions

#### carts
User shopping carts.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Cart identifier |
| user_id | UUID | UNIQUE, NOT NULL | Cart owner |
| created_at | TIMESTAMP | NOT NULL | Cart creation time |
| updated_at | TIMESTAMP | NOT NULL | Last modification time |
| expires_at | TIMESTAMP | NOT NULL | Cart expiration time |

#### cart_items
Items in shopping carts.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Cart item identifier |
| cart_id | UUID | FOREIGN KEY (carts.id), NOT NULL | Associated cart |
| product_id | UUID | NOT NULL | Product reference |
| quantity | INT | NOT NULL, CHECK (quantity > 0) | Item quantity |
| price | DECIMAL(10,2) | NOT NULL | Price snapshot at add time |
| created_at | TIMESTAMP | NOT NULL | Item added time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |

---

## Wishlist Service Database

**Database Name:** `wishlist_db`

### ER Diagram

```mermaid
erDiagram
    wishlists ||--o{ wishlist_items : contains
    
    wishlists {
        uuid id PK
        uuid user_id UK
        timestamp created_at
        timestamp updated_at
    }
    
    wishlist_items {
        uuid id PK
        uuid wishlist_id FK
        uuid product_id
        timestamp added_at
        timestamp created_at
    }
```

### Table Descriptions

#### wishlists
User wishlists.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Wishlist identifier |
| user_id | UUID | UNIQUE, NOT NULL | Wishlist owner |
| created_at | TIMESTAMP | NOT NULL | Wishlist creation time |
| updated_at | TIMESTAMP | NOT NULL | Last modification time |

#### wishlist_items
Items in wishlists.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Wishlist item identifier |
| wishlist_id | UUID | FOREIGN KEY (wishlists.id), NOT NULL | Associated wishlist |
| product_id | UUID | NOT NULL | Product reference |
| added_at | TIMESTAMP | NOT NULL | When item was added |
| created_at | TIMESTAMP | NOT NULL | Record creation time |

---

## Order Service Database

**Database Name:** `order_db`

### ER Diagram

```mermaid
erDiagram
    orders ||--o{ order_items : contains
    orders ||--o{ order_status_history : tracks
    
    orders {
        uuid id PK
        uuid user_id
        varchar order_number UK
        varchar status
        decimal subtotal
        decimal tax
        decimal shipping_cost
        decimal total
        jsonb shipping_address
        jsonb billing_address
        timestamp created_at
        timestamp updated_at
    }
    
    order_items {
        uuid id PK
        uuid order_id FK
        uuid product_id
        varchar product_name
        decimal product_price
        int quantity
        decimal subtotal
        timestamp created_at
    }
    
    order_status_history {
        uuid id PK
        uuid order_id FK
        varchar status
        text comment
        timestamp created_at
    }
```

### Table Descriptions

#### orders
Customer orders.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Order identifier |
| user_id | UUID | NOT NULL | Customer user ID |
| order_number | VARCHAR(50) | UNIQUE, NOT NULL | Human-readable order number |
| status | VARCHAR(20) | NOT NULL | Order status (pending, confirmed, etc.) |
| subtotal | DECIMAL(10,2) | NOT NULL | Items subtotal |
| tax | DECIMAL(10,2) | NOT NULL | Tax amount |
| shipping_cost | DECIMAL(10,2) | NOT NULL | Shipping cost |
| total | DECIMAL(10,2) | NOT NULL | Total order amount |
| shipping_address | JSONB | NOT NULL | Shipping address details |
| billing_address | JSONB | NOT NULL | Billing address details |
| created_at | TIMESTAMP | NOT NULL | Order creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |

#### order_items
Items in orders (product snapshots).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Order item identifier |
| order_id | UUID | FOREIGN KEY (orders.id), NOT NULL | Associated order |
| product_id | UUID | NOT NULL | Product reference |
| product_name | VARCHAR(255) | NOT NULL | Product name snapshot |
| product_price | DECIMAL(10,2) | NOT NULL | Product price snapshot |
| quantity | INT | NOT NULL, CHECK (quantity > 0) | Item quantity |
| subtotal | DECIMAL(10,2) | NOT NULL | Line item total |
| created_at | TIMESTAMP | NOT NULL | Record creation time |

#### order_status_history
Audit trail of order status changes.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | History record identifier |
| order_id | UUID | FOREIGN KEY (orders.id), NOT NULL | Associated order |
| status | VARCHAR(20) | NOT NULL | New status |
| comment | TEXT | | Status change notes |
| created_at | TIMESTAMP | NOT NULL | When status changed |

---

## Payment Service Database

**Database Name:** `payment_db`

### ER Diagram

```mermaid
erDiagram
    payments ||--o{ refunds : has
    
    payments {
        uuid id PK
        uuid order_id UK
        uuid user_id
        varchar stripe_payment_intent_id UK
        decimal amount
        varchar currency
        varchar status
        varchar payment_method
        timestamp created_at
        timestamp updated_at
    }
    
    refunds {
        uuid id PK
        uuid payment_id FK
        varchar stripe_refund_id UK
        decimal amount
        text reason
        varchar status
        timestamp created_at
    }
```

### Table Descriptions

#### payments
Payment transactions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Payment identifier |
| order_id | UUID | UNIQUE, NOT NULL | Associated order |
| user_id | UUID | NOT NULL | Customer user ID |
| stripe_payment_intent_id | VARCHAR(255) | UNIQUE | Stripe payment intent ID |
| amount | DECIMAL(10,2) | NOT NULL | Payment amount |
| currency | VARCHAR(3) | DEFAULT 'USD' | Currency code |
| status | VARCHAR(20) | NOT NULL | Payment status |
| payment_method | VARCHAR(50) | | Payment method used |
| created_at | TIMESTAMP | NOT NULL | Payment creation time |
| updated_at | TIMESTAMP | NOT NULL | Last update time |

#### refunds
Payment refunds.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Refund identifier |
| payment_id | UUID | FOREIGN KEY (payments.id), NOT NULL | Associated payment |
| stripe_refund_id | VARCHAR(255) | UNIQUE | Stripe refund ID |
| amount | DECIMAL(10,2) | NOT NULL | Refund amount |
| reason | TEXT | | Refund reason |
| status | VARCHAR(20) | NOT NULL | Refund status |
| created_at | TIMESTAMP | NOT NULL | Refund creation time |

---

## Inventory Service Database

**Database Name:** `inventory_db`

### ER Diagram

```mermaid
erDiagram
    warehouses ||--o{ inventory : stores
    inventory ||--o{ inventory_transactions : tracks
    
    warehouses {
        uuid id PK
        varchar name
        varchar code UK
        jsonb address
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }
    
    inventory {
        uuid id PK
        uuid product_id UK
        varchar sku
        int quantity_available
        int quantity_reserved
        int quantity_sold
        uuid warehouse_id FK
        int low_stock_threshold
        timestamp created_at
        timestamp updated_at
    }
    
    inventory_transactions {
        uuid id PK
        uuid inventory_id FK
        varchar transaction_type
        int quantity
        int previous_qty
        int new_qty
        uuid reference_id
        text notes
        timestamp created_at
    }
```

### Table Descriptions

#### warehouses
Warehouse locations.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Warehouse identifier |
| name | VARCHAR(255) | NOT NULL | Warehouse name |
| code | VARCHAR(50) | UNIQUE, NOT NULL | Warehouse code |
| address | JSONB | NOT NULL | Warehouse address |
| is_active | BOOLEAN | DEFAULT true | Active status |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

#### inventory
Product inventory levels.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Inventory record identifier |
| product_id | UUID | UNIQUE, NOT NULL | Product reference |
| sku | VARCHAR(100) | NOT NULL | Product SKU |
| quantity_available | INT | NOT NULL, DEFAULT 0 | Available stock |
| quantity_reserved | INT | NOT NULL, DEFAULT 0 | Reserved stock |
| quantity_sold | INT | NOT NULL, DEFAULT 0 | Total sold |
| warehouse_id | UUID | FOREIGN KEY (warehouses.id) | Warehouse location |
| low_stock_threshold | INT | DEFAULT 10 | Low stock alert threshold |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |

#### inventory_transactions
Audit trail of inventory changes.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Transaction identifier |
| inventory_id | UUID | FOREIGN KEY (inventory.id), NOT NULL | Associated inventory |
| transaction_type | VARCHAR(20) | NOT NULL | Type (purchase, sale, adjustment, return) |
| quantity | INT | NOT NULL | Quantity changed |
| previous_qty | INT | NOT NULL | Quantity before change |
| new_qty | INT | NOT NULL | Quantity after change |
| reference_id | UUID | | Reference to order, etc. |
| notes | TEXT | | Transaction notes |
| created_at | TIMESTAMP | NOT NULL | Transaction timestamp |

---

## Notification Service Database

**Database Name:** `notification_db`

### ER Diagram

```mermaid
erDiagram
    notifications {
        uuid id PK
        uuid user_id
        varchar type
        varchar template
        varchar recipient
        varchar subject
        text body
        varchar status
        int retry_count
        text error_message
        timestamp sent_at
        timestamp created_at
    }
```

### Table Descriptions

#### notifications
Notification history and status.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Notification identifier |
| user_id | UUID | | Recipient user ID |
| type | VARCHAR(50) | NOT NULL | Notification type (email, sms, push) |
| template | VARCHAR(100) | NOT NULL | Template name used |
| recipient | VARCHAR(255) | NOT NULL | Recipient address |
| subject | VARCHAR(255) | NOT NULL | Email subject |
| body | TEXT | NOT NULL | Email body (HTML) |
| status | VARCHAR(20) | DEFAULT 'pending' | Status (pending, sent, failed) |
| retry_count | INT | DEFAULT 0 | Number of retry attempts |
| error_message | TEXT | | Error message if failed |
| sent_at | TIMESTAMP | | When successfully sent |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |

---

## Reporting Service Database

**Database Name:** `reporting_db`

### ER Diagram

```mermaid
erDiagram
    daily_metrics {
        uuid id PK
        date date UK
        decimal revenue
        int orders_count
        int new_customers
        jsonb metrics_json
        timestamp created_at
    }
    
    weekly_metrics {
        uuid id PK
        int year
        int week
        decimal revenue
        int orders_count
        int new_customers
        jsonb metrics_json
        timestamp created_at
    }
    
    monthly_metrics {
        uuid id PK
        int year
        int month
        decimal revenue
        int orders_count
        int new_customers
        jsonb metrics_json
        timestamp created_at
    }
```

### Table Descriptions

#### daily_metrics
Pre-aggregated daily metrics.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Metric identifier |
| date | DATE | UNIQUE, NOT NULL | Metric date |
| revenue | DECIMAL(10,2) | NOT NULL | Total revenue |
| orders_count | INT | NOT NULL | Total orders |
| new_customers | INT | NOT NULL | New customer count |
| metrics_json | JSONB | | Additional metrics |
| created_at | TIMESTAMP | NOT NULL | Record creation time |

#### weekly_metrics
Pre-aggregated weekly metrics.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Metric identifier |
| year | INT | NOT NULL | Year |
| week | INT | NOT NULL | Week number |
| revenue | DECIMAL(10,2) | NOT NULL | Total revenue |
| orders_count | INT | NOT NULL | Total orders |
| new_customers | INT | NOT NULL | New customer count |
| metrics_json | JSONB | | Additional metrics |
| created_at | TIMESTAMP | NOT NULL | Record creation time |

Unique constraint: (year, week)

#### monthly_metrics
Pre-aggregated monthly metrics.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Metric identifier |
| year | INT | NOT NULL | Year |
| month | INT | NOT NULL | Month (1-12) |
| revenue | DECIMAL(10,2) | NOT NULL | Total revenue |
| orders_count | INT | NOT NULL | Total orders |
| new_customers | INT | NOT NULL | New customer count |
| metrics_json | JSONB | | Additional metrics |
| created_at | TIMESTAMP | NOT NULL | Record creation time |

Unique constraint: (year, month)

---

## Database Indexes

### Auth Service Indexes

```sql
-- users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);

-- refresh_tokens table
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- password_resets table
CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);
CREATE INDEX idx_password_resets_token ON password_resets(token);
```

### Product Service Indexes

```sql
-- products table
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_average_rating ON products(average_rating);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_name_trgm ON products USING gin(name gin_trgm_ops); -- Full-text search

-- categories table
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug ON categories(slug);

-- reviews table
CREATE INDEX idx_reviews_product_id ON reviews(product_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);
CREATE INDEX idx_reviews_created_at ON reviews(created_at);
```

### Cart Service Indexes

```sql
-- carts table
CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_expires_at ON carts(expires_at);

-- cart_items table
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);
```

### Order Service Indexes

```sql
-- orders table
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_order_number ON orders(order_number);

-- order_items table
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- order_status_history table
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_order_status_history_created_at ON order_status_history(created_at);
```

### Payment Service Indexes

```sql
-- payments table
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);
CREATE INDEX idx_payments_stripe_payment_intent_id ON payments(stripe_payment_intent_id);

-- refunds table
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_status ON refunds(status);
```

### Inventory Service Indexes

```sql
-- inventory table
CREATE INDEX idx_inventory_product_id ON inventory(product_id);
CREATE INDEX idx_inventory_warehouse_id ON inventory(warehouse_id);
CREATE INDEX idx_inventory_sku ON inventory(sku);
CREATE INDEX idx_inventory_quantity_available ON inventory(quantity_available);

-- inventory_transactions table
CREATE INDEX idx_inventory_transactions_inventory_id ON inventory_transactions(inventory_id);
CREATE INDEX idx_inventory_transactions_reference_id ON inventory_transactions(reference_id);
CREATE INDEX idx_inventory_transactions_created_at ON inventory_transactions(created_at);
```

### Notification Service Indexes

```sql
-- notifications table
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
```

### Reporting Service Indexes

```sql
-- daily_metrics table
CREATE INDEX idx_daily_metrics_date ON daily_metrics(date);

-- weekly_metrics table
CREATE INDEX idx_weekly_metrics_year_week ON weekly_metrics(year, week);

-- monthly_metrics table
CREATE INDEX idx_monthly_metrics_year_month ON monthly_metrics(year, month);
```

---

## Migration Strategy

### Tool
**golang-migrate** - Database migration tool for Go

### Migration Files
Location: `{service}/migrations/`

Naming convention: `{version}_{description}.{up|down}.sql`

Example:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

### Migration Execution

```bash
# Run migrations
migrate -path ./migrations -database "postgresql://user:pass@localhost:5432/auth_db?sslmode=disable" up

# Rollback one migration
migrate -path ./migrations -database "postgresql://user:pass@localhost:5432/auth_db?sslmode=disable" down 1

# Check migration version
migrate -path ./migrations -database "postgresql://user:pass@localhost:5432/auth_db?sslmode=disable" version
```

### Example Migration: Auth Service

**000001_create_users_table.up.sql:**
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(20) DEFAULT 'customer' CHECK (role IN ('customer', 'admin')),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

**000001_create_users_table.down.sql:**
```sql
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS users;
```

### Seed Data

Location: `{service}/migrations/seeds/`

Example seed file for products:
```sql
-- Insert categories
INSERT INTO categories (id, name, slug, description) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Electronics', 'electronics', 'Electronic devices and gadgets'),
    ('550e8400-e29b-41d4-a716-446655440002', 'Clothing', 'clothing', 'Apparel and fashion'),
    ('550e8400-e29b-41d4-a716-446655440003', 'Books', 'books', 'Books and literature');

-- Insert products
INSERT INTO products (id, name, description, price, category_id, brand, sku) VALUES
    ('550e8400-e29b-41d4-a716-446655440101', 'Laptop', 'High-performance laptop', 999.99, '550e8400-e29b-41d4-a716-446655440001', 'TechBrand', 'LAPTOP-001'),
    ('550e8400-e29b-41d4-a716-446655440102', 'Smartphone', 'Latest smartphone model', 699.99, '550e8400-e29b-41d4-a716-446655440001', 'PhoneCo', 'PHONE-001');
```

---

## Database Best Practices

### 1. Connection Pooling
```go
// Configure connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 2. Prepared Statements
Use parameterized queries to prevent SQL injection and improve performance.

### 3. Transactions
Use transactions for operations that modify multiple tables.

### 4. Soft Deletes
Use `deleted_at` timestamp for soft deletes to maintain data integrity.

### 5. UUID Primary Keys
Use UUIDs for better distribution and no central point of failure.

### 6. JSONB for Flexibility
Use JSONB for semi-structured data like addresses, metadata.

### 7. Regular Backups
- Automated daily backups
- Point-in-time recovery (PITR)
- Backup retention: 30 days

### 8. Monitoring
- Connection pool metrics
- Query performance
- Slow query logs
- Database size growth

---

This database design follows normalization principles, provides proper indexing for performance, and maintains data integrity through foreign keys and constraints.
