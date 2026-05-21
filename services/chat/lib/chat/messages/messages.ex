defmodule Chat.Messages do
  import Ecto.Query

  alias Chat.Messages.Message
  alias Chat.Repo

  def list_recent(guild_id, limit \\ 50) do
    Message
    |> where([m], m.guild_id == ^guild_id)
    |> order_by([m], desc: m.created_at, desc: m.id)
    |> limit(^limit)
    |> Repo.all()
    |> Enum.reverse()
  end

  def create(guild_id, user_id, body) do
    id = "chat_message_" <> Uniq.UUID.uuid7(:hex)

    %Message{}
    |> Message.changeset(%{id: id, guild_id: guild_id, user_id: user_id, body: body})
    |> Repo.insert()
  end

  def format_message(%Message{} = msg) do
    %{
      id: msg.id,
      guild_id: msg.guild_id,
      user_id: msg.user_id,
      body: msg.body,
      created_at: DateTime.to_iso8601(msg.created_at)
    }
  end
end
