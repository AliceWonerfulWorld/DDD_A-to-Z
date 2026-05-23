FROM node:22-alpine AS node-builder
WORKDIR /app
COPY package*.json vite.config.ts ./
RUN npm ci
COPY resources ./resources
RUN npm run build

FROM composer:2 AS php-builder
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --optimize-autoloader --no-interaction

FROM dunglas/frankenphp AS runtime
WORKDIR /app

RUN install-php-extensions pdo pdo_pgsql

COPY . .
COPY --from=node-builder /app/public/build ./public/build
COPY --from=php-builder /app/vendor ./vendor
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh \
  && groupadd -r admin \
  && useradd -r -g admin admin \
  && chown -R admin:admin /app

USER admin

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["frankenphp", "php-server", "--listen", "0.0.0.0:8080", "--document-root", "/app/public"]
