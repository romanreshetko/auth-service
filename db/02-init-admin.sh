#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
INSERT INTO users (email, nickname, password, user_role, email_verified)
SELECT '$ADMIN_EMAIL', 'admin', '$ADMIN_HASH', 'admin', true
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE user_role = 'admin'
);
EOSQL