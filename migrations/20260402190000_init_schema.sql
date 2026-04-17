-- +goose Up
CREATE TABLE positions (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    size_liters REAL NOT NULL CHECK (size_liters > 0),
    quantity INT NOT NULL CHECK (quantity >= 0),
    price BIGINT NOT NULL CHECK (price > 0),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE clients (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE sellers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('new', 'ready', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_orders_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE RESTRICT,
    CONSTRAINT fk_orders_seller FOREIGN KEY (seller_id) REFERENCES sellers(id) ON DELETE RESTRICT
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    position_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price BIGINT NOT NULL CHECK (price > 0),
    CONSTRAINT uq_order_items_order_position UNIQUE (order_id, position_id),
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE RESTRICT
);

CREATE INDEX idx_orders_client_id ON orders(client_id);
CREATE INDEX idx_orders_seller_id ON orders(seller_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_position_id ON order_items(position_id);

-- +goose Down
DROP INDEX IF EXISTS idx_order_items_position_id;
DROP INDEX IF EXISTS idx_order_items_order_id;
DROP INDEX IF EXISTS idx_orders_seller_id;
DROP INDEX IF EXISTS idx_orders_client_id;

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS sellers;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS positions;
