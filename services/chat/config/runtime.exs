import Config

if config_env() == :prod do
  config :chat, Chat.Repo,
    url: System.fetch_env!("DATABASE_URL"),
    pool_size: String.to_integer(System.get_env("POOL_SIZE") || "5"),
    ssl: true,
    ssl_opts: [verify: :verify_none]

  config :chat, Chat.Endpoint,
    http: [ip: {0, 0, 0, 0}, port: String.to_integer(System.get_env("PORT") || "4000")],
    server: true,
    secret_key_base: System.fetch_env!("SECRET_KEY_BASE"),
    check_origin: [System.fetch_env!("ALLOWED_ORIGIN")]
end
