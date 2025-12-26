CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    total_qty INT NOT NULL DEFAULT 0,
    reserved_qty INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);

INSERT INTO products (sku, name, description, price, total_qty, reserved_qty) VALUES
('PROD-001', 'High-Performance Laptop', 'A powerful laptop for developers.', 1999.99, 100, 0),
('PROD-002', 'Wireless Noise-Canceling Headphones', 'Immersive sound experience.', 299.99, 200, 0),
('PROD-003', 'Smartphone X', 'Latest generation smartphone with advanced camera.', 999.99, 150, 0);

