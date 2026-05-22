defmodule Chat.Router do
  use Phoenix.Router, helpers: false

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/", Chat do
    pipe_through :api
    get "/health", HealthController, :index
  end
end
