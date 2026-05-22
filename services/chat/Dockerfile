FROM hexpm/elixir:1.17-erlang-27-debian-bullseye-20250317-slim AS build

WORKDIR /app

RUN mix local.hex --force && mix local.rebar --force

COPY mix.exs mix.lock ./
RUN MIX_ENV=prod mix deps.get --only prod

COPY . .
RUN MIX_ENV=prod mix compile && mix release

FROM debian:bullseye-slim AS runtime

RUN apt-get update && \
    apt-get install -y libstdc++6 openssl libncurses5 locales && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/_build/prod/rel/chat ./

RUN groupadd -r chat && useradd -r -g chat chat && \
    chown -R chat:chat /app

USER chat

CMD ["./bin/chat", "start"]
