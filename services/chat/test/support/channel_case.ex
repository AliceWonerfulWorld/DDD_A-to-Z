defmodule Chat.ChannelCase do
  use ExUnit.CaseTemplate

  using do
    quote do
      import Phoenix.ChannelTest
      import Chat.Factory

      @endpoint Chat.Endpoint
    end
  end

  setup tags do
    Chat.DataCase.setup_sandbox(tags)
  end
end
