import Config

config :chat, Chat.Repo,
  url: System.get_env("DATABASE_URL") || "postgres://lang_war:1234567890@localhost:5432/lang_war",
  pool_size: 5

config :chat, Chat.Endpoint,
  http: [ip: {127, 0, 0, 1}, port: 4000],
  server: true,
  check_origin: ["http://localhost:5173"]

config :logger, level: :debug
