-- Создаем пользователей
CREATE USER sentry WITH PASSWORD 'sentry';
CREATE USER keycloak WITH PASSWORD 'keycloak';
-- Создаем базы данных
CREATE DATABASE sentry OWNER sentry;
CREATE DATABASE keycloak OWNER keycloak;
-- Даем права
ALTER USER sentry WITH SUPERUSER;
ALTER USER "chat-service" WITH CREATEDB SUPERUSER;
GRANT ALL PRIVILEGES ON DATABASE sentry TO sentry;
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;