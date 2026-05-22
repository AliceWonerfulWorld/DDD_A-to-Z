defmodule Chat.DataCase do
  use ExUnit.CaseTemplate

  using do
    quote do
      import Chat.Factory
    end
  end

  def setup_sandbox(tags) do
    pid = Ecto.Adapters.SQL.Sandbox.start_owner!(Chat.Repo, shared: not tags[:async])
    on_exit(fn -> Ecto.Adapters.SQL.Sandbox.stop_owner(pid) end)
    :ok
  end

  setup tags do
    setup_sandbox(tags)
  end
end
