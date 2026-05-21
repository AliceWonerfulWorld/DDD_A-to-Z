defmodule Chat.Messages.Message do
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :string, autogenerate: false}
  @timestamps_opts [inserted_at: :created_at, updated_at: false, type: :utc_datetime_usec]

  schema "guild_chat_messages" do
    field :guild_id, :string
    field :user_id, :string
    field :body, :string
    timestamps()
  end

  def changeset(msg, attrs) do
    msg
    |> cast(attrs, [:id, :guild_id, :user_id, :body])
    |> validate_required([:id, :guild_id, :user_id, :body])
    |> validate_length(:body, min: 1, max: 1000)
  end
end
