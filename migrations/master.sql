CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscription_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    duration INTERVAL NOT NULL,
    max_devices INT NOT NULL
);

INSERT INTO subscription_types (name, duration, max_devices)
VALUES 
('basic', INTERVAL '30 days', 3),
('premium', INTERVAL '90 days', 10);

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    subscription_id INT NOT NULL REFERENCES subscription_types(id),
    expiry_date TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS action_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

INSERT INTO action_types (name)
VALUES 
    ('create'),
    ('delete');

CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    action_id INT NOT NULL REFERENCES action_types(id),
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS servers (
  id SERIAL PRIMARY KEY,
  protocol TEXT UNIQUE NOT NULL,
  ip_address INET UNIQUE NOT NULL,
  port INT NOT NULL
);