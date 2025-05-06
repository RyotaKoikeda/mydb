BEGIN;

INSERT INTO dummy_orders (user_id, product_name, price)
VALUES 
  (1, 'Book', 1200),
  (2, 'Pen', 200),
  (1, 'Notebook', 500);

COMMIT;
