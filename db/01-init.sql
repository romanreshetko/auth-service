CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    nickname TEXT NOT NULL,
    password TEXT NOT NULL,
    user_role TEXT NOT NULL CHECK (user_role IN ('user', 'moderator', 'admin')),
    photo TEXT,
    city BIGINT,
    status TEXT,
    points INTEGER,
    agreement_pd BOOlEAN,
    agreement_ea BOOlEAN,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS email_verifications (
    id BIGSERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    code TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NOW() + interval '10 minutes'
);

CREATE OR REPLACE FUNCTION update_user_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.points >= 401 THEN
       NEW.status = 'Вечный странник';
    ELSIF NEW.points >= 201 THEN
        NEW.status = 'Легенда дорог';
    ELSIF NEW.points >= 101 THEN
        NEW.status = 'Пилигрим';
    ELSIF NEW.points >= 51 THEN
        NEW.status = 'Исследователь';
    ELSE
        NEW.status = 'Новичок';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_status
BEFORE UPDATE OF points ON users
FOR EACH ROW
EXECUTE FUNCTION update_user_status();

INSERT INTO users (email, nickname, password, user_role, city, status, points, agreement_pd, agreement_ea, email_verified)
VALUES ('test2@mail.com', 'Алекс', '$2a$10$SwddE7VLxCNI8L8bD5j6TuqZ4HMHrdZjADAaZfbaXvmBWCL0IhEzO',
        'user', 1, 'Новичок', 10, true, true, true);