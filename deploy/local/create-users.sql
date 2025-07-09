-- Создаем пользователей только если они не существуют
DO $$ BEGIN IF NOT EXISTS (
  SELECT
  FROM pg_catalog.pg_roles
  WHERE rolname = 'sentry'
) THEN CREATE ROLE sentry WITH LOGIN PASSWORD 'sentry';
GRANT ALL PRIVILEGES ON DATABASE sentry TO sentry;
ALTER USER sentry WITH SUPERUSER;
END IF;
IF NOT EXISTS (
  SELECT
  FROM pg_catalog.pg_roles
  WHERE rolname = 'keycloak'
) THEN CREATE ROLE keycloak WITH LOGIN PASSWORD 'keycloak';
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;
END IF;
-- Даем права пользователю chat-service создавать БД для тестов
ALTER USER "chat-service" WITH CREATEDB SUPERUSER;
END $$;
-- Создаем базы данных только если они не существуют
SELECT 'CREATE DATABASE sentry'
WHERE NOT EXISTS (
    SELECT
    FROM pg_database
    WHERE datname = 'sentry'
  ) \ gexec
SELECT 'CREATE DATABASE keycloak'
WHERE NOT EXISTS (
    SELECT
    FROM pg_database
    WHERE datname = 'keycloak'
  ) \ gexec