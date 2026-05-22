import Config

if config_env() == :prod do
  database_url = System.fetch_env!("DATABASE_URL")
  database_host = URI.parse(database_url).host |> to_charlist()

  ssl_opts =
    case System.get_env("DB_CA_CERT_PATH") do
      nil -> [cacerts: :public_key.cacerts_get()]
      path -> [cacertfile: path]
    end
    |> Keyword.merge(
      verify: :verify_peer,
      server_name_indication: database_host,
      customize_hostname_check: [
        match_fun: :public_key.pkix_verify_hostname_match_fun(:https)
      ]
    )

  config :chat, Chat.Repo,
    url: database_url,
    pool_size: String.to_integer(System.get_env("POOL_SIZE") || "5"),
    ssl: ssl_opts

  config :chat, Chat.Endpoint,
    http: [ip: {0, 0, 0, 0}, port: String.to_integer(System.get_env("PORT") || "4000")],
    server: true,
    secret_key_base: System.fetch_env!("SECRET_KEY_BASE"),
    check_origin: [System.fetch_env!("ALLOWED_ORIGIN")]
end
