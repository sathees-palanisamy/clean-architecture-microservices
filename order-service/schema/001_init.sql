CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    order_status VARCHAR(50) NOT NULL,
    payment_status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

INSERT INTO orders (user_id, product_id, product_name, unit_price, quantity, total_price, order_status, payment_status) VALUES
(101, 1, 'High-Performance Laptop', 1999.99, 1, 1999.99, 'COMPLETED', 'PAID'),
(102, 2, 'Wireless Noise-Canceling Headphones', 299.99, 2, 599.98, 'PENDING', 'PENDING');

