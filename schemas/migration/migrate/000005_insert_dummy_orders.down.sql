BEGIN;

DELETE FROM dummy_orders
WHERE product_name IN ('Book', 'Pen', 'Notebook');

COMMIT;
