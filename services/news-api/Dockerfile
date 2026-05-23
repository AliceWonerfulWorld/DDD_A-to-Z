FROM node:22-slim AS build

WORKDIR /app

COPY package.json ./
RUN npm install --ignore-scripts

COPY tsconfig.json ./
COPY src/ ./src/
RUN npm run build

FROM node:22-slim AS runtime

WORKDIR /app

COPY package.json ./
RUN npm install --omit=dev --ignore-scripts

COPY --from=build /app/dist ./dist

RUN groupadd -r app && useradd -r -g app app && \
    chown -R app:app /app

USER app

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD ["node", "-e", "fetch('http://localhost:8080/healthz').then(r => { process.exit(r.ok ? 0 : 1) }).catch(() => process.exit(1))"]

EXPOSE 8080

CMD ["node", "dist/index.js"]
