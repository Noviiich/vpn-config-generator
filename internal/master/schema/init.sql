CREATE TABLE IF NOT EXISTS users (
    telegram_id BIGINT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    subscription_active BOOLEAN DEFAULT FALSE,
    subscription_expiry TIMESTAMP
);

CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    country_flag TEXT,
    ip_address TEXT NOT NULL,
    port INTEGER NOT NULL,
    protocol TEXT NOT NULL
);