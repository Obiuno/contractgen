BEGIN;

-- ==========================================
-- TABLES
-- ==========================================

-- Table: customers
CREATE TABLE IF NOT EXISTS customers (
	id UUID,
	email VARCHAR,
	name VARCHAR,
	phone VARCHAR,
	created_at TIMESTAMP
);

-- Table: products
CREATE TABLE IF NOT EXISTS products (
	id UUID,
	sku VARCHAR,
	name VARCHAR,
	price DECIMAL
);

-- Table: orders
CREATE TABLE IF NOT EXISTS orders (
	id UUID,
	customer_id UUID,
	order_total DECIMAL,
	status VARCHAR,
	created_at TIMESTAMP
);

-- Table: order_items
CREATE TABLE IF NOT EXISTS order_items (
	id UUID,
	order_id UUID,
	product_id UUID,
	quantity INTEGER,
	unit_price DECIMAL
);

-- ==========================================
-- CONSTRAINTS (PK, FK, UNIQUE, NOT NULL)
-- ==========================================

-- Constraints for customers
ALTER TABLE customers ALTER COLUMN id SET NOT NULL;
ALTER TABLE customers ADD CONSTRAINT pk_customers PRIMARY KEY (id);
ALTER TABLE customers ALTER COLUMN email SET NOT NULL;
ALTER TABLE customers ADD CONSTRAINT uq_customers_email UNIQUE (email);
ALTER TABLE customers ALTER COLUMN name SET NOT NULL;
ALTER TABLE customers ALTER COLUMN created_at SET NOT NULL;

-- Constraints for products
ALTER TABLE products ALTER COLUMN id SET NOT NULL;
ALTER TABLE products ADD CONSTRAINT pk_products PRIMARY KEY (id);
ALTER TABLE products ALTER COLUMN sku SET NOT NULL;
ALTER TABLE products ADD CONSTRAINT uq_products_sku UNIQUE (sku);
ALTER TABLE products ALTER COLUMN name SET NOT NULL;
ALTER TABLE products ALTER COLUMN price SET NOT NULL;

-- Constraints for orders
ALTER TABLE orders ALTER COLUMN id SET NOT NULL;
ALTER TABLE orders ADD CONSTRAINT pk_orders PRIMARY KEY (id);
ALTER TABLE orders ALTER COLUMN customer_id SET NOT NULL;
ALTER TABLE orders ADD CONSTRAINT fk_orders_customer_id FOREIGN KEY (customer_id) REFERENCES customers (id);
ALTER TABLE orders ALTER COLUMN order_total SET NOT NULL;
ALTER TABLE orders ALTER COLUMN status SET NOT NULL;
ALTER TABLE orders ALTER COLUMN created_at SET NOT NULL;

-- Constraints for order_items
ALTER TABLE order_items ALTER COLUMN id SET NOT NULL;
ALTER TABLE order_items ADD CONSTRAINT pk_order_items PRIMARY KEY (id);
ALTER TABLE order_items ALTER COLUMN order_id SET NOT NULL;
ALTER TABLE order_items ADD CONSTRAINT fk_order_items_order_id FOREIGN KEY (order_id) REFERENCES orders (id);
ALTER TABLE order_items ALTER COLUMN product_id SET NOT NULL;
ALTER TABLE order_items ADD CONSTRAINT fk_order_items_product_id FOREIGN KEY (product_id) REFERENCES products (id);
ALTER TABLE order_items ALTER COLUMN quantity SET NOT NULL;
ALTER TABLE order_items ALTER COLUMN unit_price SET NOT NULL;

COMMIT;
