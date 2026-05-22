defmodule Chat.GuildChannel do
  use Phoenix.Channel

  import Ecto.Query, only: [from: 2]

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
    query =
      from(m in "guild_memberships",
        where:
          m.user_id == ^user_id and
            m.guild_id == ^guild_id and
            is_nil(m.left_at),
        select: 1,
        limit: 1
      )

    Repo.exists?(query)
  end
end
