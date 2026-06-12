-- Drop triggers
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
