FROM node:22-alpine AS node-builder
WORKDIR /workspace
RUN npm install -g pnpm
# ルートの lockfile と admin の package.json をコピー
COPY pnpm-lock.yaml pnpm-workspace.yaml package.json ./
COPY services/admin/package.json ./services/admin/
RUN pnpm install --frozen-lockfile --filter admin
COPY services/admin/vite.config.ts services/admin/tsconfig*.json ./services/admin/
COPY services/admin/resources ./services/admin/resources
WORKDIR /workspace/services/admin
RUN pnpm run build

FROM composer:2 AS php-builder
WORKDIR /app
COPY composer.json composer.lock ./
RUN composer install --no-dev --optimize-autoloader --no-interaction

FROM dunglas/frankenphp AS runtime
WORKDIR /app

RUN install-php-extensions pdo pdo_pgsql

COPY . .
COPY --from=node-builder /workspace/services/admin/public/build ./public/build
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
