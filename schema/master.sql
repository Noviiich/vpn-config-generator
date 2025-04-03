CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS subscription_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    duration INTERVAL NOT NULL,
    max_devices INT NOT NULL
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    type_id INT NOT NULL REFERENCES subscription_types(id),
    expiry_date TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    action INT NOT NULL,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);