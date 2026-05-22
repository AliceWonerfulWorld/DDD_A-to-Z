import Config

config :chat, Chat.Repo,
  adapter: Ecto.Adapters.Postgres

config :chat, Chat.Endpoint,
  pubsub_server: Chat.PubSub,
  render_errors: [formats: [json: Chat.ErrorJSON]]

config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

import_config "#{config_env()}.exs"
