defmodule Chat.UserSocket do
  use Phoenix.Socket

  channel "guild:*", Chat.GuildChannel

  @impl true
  def connect(%{"token" => token}, socket, _connect_info) do
    token_hash = hash(token)

    case consume_token(token_hash) do
      {:ok, %{user_id: user_id, guild_id: guild_id}} ->
        {:ok, assign(socket, user_id: user_id, guild_id: guild_id)}

      :error ->
        :error
    end
  end

  def connect(_params, _socket, _connect_info), do: :error

  @impl true
  def id(socket), do: "users_socket:#{socket.assigns.user_id}"

  defp hash(token) do
    :crypto.hash(:sha256, token) |> Base.encode16(case: :lower)
  end

  defp consume_token(token_hash) do
    now = DateTime.utc_now()

    case Chat.Repo.query(
           """
           UPDATE chat_tokens
           SET used_at = NOW()
           WHERE token_hash = $1
             AND expires_at > $2
             AND used_at IS NULL
           RETURNING user_id, guild_id
           """,
           [token_hash, now]
         ) do
      {:ok, %{rows: [[user_id, guild_id]]}} ->
        {:ok, %{user_id: user_id, guild_id: guild_id}}

      {:ok, %{rows: []}} ->
        :error

      {:error, _} ->
        :error
    end
  end
end
