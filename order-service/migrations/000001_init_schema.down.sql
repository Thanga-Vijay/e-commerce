-- Drop triggers
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
