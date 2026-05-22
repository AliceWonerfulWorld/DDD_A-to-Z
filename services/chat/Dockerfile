FROM hexpm/elixir:1.17.3-erlang-27.2.4-debian-bullseye-20260518-slim AS build

WORKDIR /app

RUN mix local.hex --force && mix local.rebar --force

COPY mix.exs mix.lock ./
RUN MIX_ENV=prod mix deps.get --only prod

COPY . .
RUN MIX_ENV=prod mix compile && MIX_ENV=prod mix release

FROM debian:bullseye-slim AS runtime

RUN apt-get update && \
    apt-get install -y ca-certificates libstdc++6 openssl libncurses5 locales && \
    rm -rf /var/lib/apt/lists/*

ENV LANG=C.UTF-8

WORKDIR /app

COPY --from=build /app/_build/prod/rel/chat ./

RUN groupadd -r chat && useradd -r -g chat chat && \
    chown -R chat:chat /app

USER chat

CMD ["./bin/chat", "start"]
