defmodule Chat.GuildChannel do
  use Phoenix.Channel

  alias Chat.Messages
  alias Chat.Repo

  @impl true
  def join("guild:" <> guild_id, _params, socket) do
    socket_guild_id = socket.assigns.guild_id
    user_id = socket.assigns.user_id

    cond do
      socket_guild_id != guild_id ->
        {:error, %{reason: "unauthorized"}}

      not active_member?(user_id, guild_id) ->
        {:error, %{reason: "unauthorized"}}

      true ->
        messages = Messages.list_recent(guild_id, 50)
        formatted = Enum.map(messages, &Messages.format_message/1)
        {:ok, %{messages: formatted}, assign(socket, :guild_id, guild_id)}
    end
  end

  @impl true
  def handle_in("new_message", %{"body" => body}, socket) do
    guild_id = socket.assigns.guild_id
    user_id = socket.assigns.user_id

    case Messages.create(guild_id, user_id, body) do
      {:ok, message} ->
        broadcast(socket, "new_message", Messages.format_message(message))
        {:reply, :ok, socket}

      {:error, _changeset} ->
        {:reply, {:error, %{reason: "failed"}}, socket}
    end
  end

  defp active_member?(user_id, guild_id) do
    case Repo.query(
           """
           SELECT 1 FROM guild_memberships
           WHERE user_id = $1
             AND guild_id = $2
             AND left_at IS NULL
           LIMIT 1
           """,
           [user_id, guild_id]
         ) do
      {:ok, %{rows: [_]}} -> true
      _ -> false
    end
  end
end
