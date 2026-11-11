-- Stage 2: database structure. Safe to run again on this development database.
-- Passwords and digital keys are NEVER returned by the public catalog API.
BEGIN;

-- Prepared for Stage 3 typo-tolerant search.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_not_blank CHECK (LENGTH(TRIM(email)) > 3)
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_unique ON users (LOWER(email));

CREATE TABLE IF NOT EXISTS products (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    platform TEXT NOT NULL,
    region TEXT NOT NULL,
    genre TEXT NOT NULL DEFAULT '',
    cover_image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT products_title_not_blank CHECK (LENGTH(TRIM(title)) > 0)
);
CREATE INDEX IF NOT EXISTS products_title_trgm_idx ON products USING GIN (title gin_trgm_ops);

CREATE TABLE IF NOT EXISTS sellers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- A product describes the game/edition/platform/region.
-- An offer describes one seller's price for that product.
CREATE TABLE IF NOT EXISTS offers (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    seller_id BIGINT NOT NULL REFERENCES sellers(id) ON DELETE RESTRICT,
    price_cents INTEGER NOT NULL CHECK (price_cents > 0),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT offers_one_per_seller_product UNIQUE (product_id, seller_id)
);
CREATE INDEX IF NOT EXISTS offers_product_id_idx ON offers (product_id);

CREATE TABLE IF NOT EXISTS orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','paid','failed')),
    total_cents INTEGER NOT NULL CHECK (total_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS orders_user_id_idx ON orders (user_id);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    offer_id BIGINT NOT NULL REFERENCES offers(id) ON DELETE RESTRICT,
    price_cents INTEGER NOT NULL CHECK (price_cents > 0)
);
CREATE INDEX IF NOT EXISTS order_items_order_id_idx ON order_items (order_id);

-- These are only fictional demo keys. Never seed real redeemable game keys.
-- At checkout, Stage 7 will allocate a free key transactionally to an order item.
CREATE TABLE IF NOT EXISTS digital_keys (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    offer_id BIGINT NOT NULL REFERENCES offers(id) ON DELETE RESTRICT,
    key_code TEXT NOT NULL UNIQUE,
    order_item_id BIGINT UNIQUE REFERENCES order_items(id) ON DELETE RESTRICT,
    delivered_at TIMESTAMPTZ,
    CONSTRAINT digital_key_delivery_consistency CHECK (
        (order_item_id IS NULL AND delivered_at IS NULL)
        OR (order_item_id IS NOT NULL AND delivered_at IS NOT NULL)
    )
);
CREATE INDEX IF NOT EXISTS digital_keys_available_idx
    ON digital_keys(offer_id, id) WHERE order_item_id IS NULL;

-- Homepage block configuration is stored in the database, never hardcoded in React.
-- Later stages will serve each block through its own Go service/endpoint.
CREATE TABLE IF NOT EXISTS homepage_sections (
    section_key TEXT PRIMARY KEY,
    heading TEXT NOT NULL,
    subtitle TEXT NOT NULL DEFAULT '',
    display_order INTEGER NOT NULL UNIQUE
);

COMMIT;
