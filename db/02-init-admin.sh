#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-'EOSQL'
INSERT INTO users (email, nickname, password, user_role, email_verified)
SELECT 'root.admin@email.com', 'admin', '$2a$10$xegv6n9MR6pUCpwomtnNbe1S0.6602IFK0YhgXSpz7BvdigwOLrBG', 'admin', true
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE user_role = 'admin'
);
EOSQL