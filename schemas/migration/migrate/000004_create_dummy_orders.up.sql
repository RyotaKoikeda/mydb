BEGIN;

CREATE TABLE dummy_orders (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    product_name VARCHAR(100),
    price INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
