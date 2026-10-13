CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT   NOT NULL,
    price_cents BIGINT NOT NULL CHECK (price_cents >= 0)
);

INSERT INTO products (name, price_cents) VALUES
    ('Keyboard', 4999),
    ('Mouse',    1999),
    ('Monitor', 18999);
