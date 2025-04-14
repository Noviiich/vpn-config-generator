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

INSERT INTO subscription_types (name, duration, max_devices)
VALUES ('basic', INTERVAL '30 days', 3);


CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    type_id INT NOT NULL REFERENCES subscription_types(id),
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
//подумай
CREATE TABLE IF NOT EXISTS device_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id INT REFERENCES device_types(id),
    private_key TEXT NOT NULL,
    public_key TEXT NOT NULL
);