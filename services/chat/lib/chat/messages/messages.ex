defmodule Chat.Messages do
  import Ecto.Query

  alias Chat.Messages.Message
  alias Chat.Repo

  def list_recent(guild_id, limit \\ 50) do
    from(m in Message,
      left_join: up in "user_profiles",
      on: up.user_id == m.user_id,
      left_join: ga in "github_accounts",
      on: ga.user_id == m.user_id,
      where: m.guild_id == ^guild_id,
      order_by: [desc: m.created_at, desc: m.id],
      limit: ^limit,
      select_merge: %{
        user_name: fragment("COALESCE(?, ?, ?)", up.display_name, ga.username, m.user_id)
      }
    )
    |> Repo.all()
    |> Enum.reverse()
  end

  def create(guild_id, user_id, body) do
    id = "chat_message_" <> Uniq.UUID.uuid7(:hex)

    %Message{}
    |> Message.changeset(%{id: id, guild_id: guild_id, user_id: user_id, body: body})
    |> Repo.insert()
    |> with_user_name()
  end

  def format_message(%Message{} = msg) do
    %{
      id: msg.id,
      guild_id: msg.guild_id,
      user_id: msg.user_id,
      user_name: msg.user_name || msg.user_id,
      body: msg.body,
      created_at: DateTime.to_iso8601(msg.created_at)
    }
  end

  defp with_user_name({:ok, %Message{} = msg}) do
    {:ok, %{msg | user_name: find_user_name(msg.user_id)}}
  end

  defp with_user_name(result), do: result

  defp find_user_name(user_id) do
    from(u in "users",
      left_join: up in "user_profiles",
      on: up.user_id == u.id,
      left_join: ga in "github_accounts",
      on: ga.user_id == u.id,
      where: u.id == ^user_id,
      select: fragment("COALESCE(?, ?, ?)", up.display_name, ga.username, u.id),
      limit: 1
    )
    |> Repo.one()
    |> Kernel.||(user_id)
  end
end
