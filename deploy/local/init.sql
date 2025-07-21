CREATE DATABASE sentry;
CREATE ROLE sentry WITH LOGIN PASSWORD 'sentry';
GRANT ALL PRIVILEGES ON DATABASE sentry TO sentry;
ALTER USER sentry WITH SUPERUSER;
-- Создаем БД keycloak
CREATE DATABASE "keycloak";
CREATE ROLE "keycloak" WITH LOGIN PASSWORD 'keycloak';
GRANT ALL PRIVILEGES ON DATABASE "keycloak" to "keycloak";
-- Даем права пользователю chat-service создавать БД для тестов
ALTER USER "chat-service" WITH CREATEDB;