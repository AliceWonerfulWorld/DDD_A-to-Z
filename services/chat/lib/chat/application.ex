defmodule Chat.Application do
  use Application

  @impl true
  def start(_type, _args) do
    children = [
      Chat.Repo,
      {Phoenix.PubSub, name: Chat.PubSub},
      Chat.Endpoint
    ]

    opts = [strategy: :one_for_one, name: Chat.Supervisor]
    Supervisor.start_link(children, opts)
  end
end
