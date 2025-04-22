CREATE TABLE IF NOT EXISTS users (
    telegram_id BIGINT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    subscription_active BOOLEAN DEFAULT FALSE,
    subscription_expiry TIMESTAMP
);