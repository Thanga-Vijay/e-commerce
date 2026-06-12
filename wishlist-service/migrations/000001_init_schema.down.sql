-- Drop triggers
DROP TRIGGER IF EXISTS update_wishlist_items_updated_at ON wishlist_items;
DROP TRIGGER IF EXISTS update_wishlists_updated_at ON wishlists;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS wishlists;
