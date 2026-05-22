defmodule Chat.Factory do
  @moduledoc false

  alias Chat.Repo

  @doc "テスト用ユーザーを直接DBに挿入する"
  def insert_user(id \\ "user_test_#{System.unique_integer([:positive])}") do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    Repo.query!(
      "INSERT INTO users (id, created_at, updated_at) VALUES ($1, $2, $2)",
      [id, now]
    )

    id
  end

  def insert_github_account(user_id, username \\ "github_#{System.unique_integer([:positive])}") do
    now = DateTime.utc_now() |> DateTime.truncate(:second)
    github_id = System.unique_integer([:positive])

    Repo.query!(
      """
      INSERT INTO github_accounts (github_id, user_id, username, avatar_url, created_at, updated_at)
      VALUES ($1, $2, $3, 'https://example.com/avatar.png', $4, $4)
      """,
      [github_id, user_id, username, now]
    )

    :ok
  end

  def insert_user_profile(user_id, display_name) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    Repo.query!(
      """
      INSERT INTO user_profiles (user_id, display_name, created_at, updated_at)
      VALUES ($1, $2, $3, $3)
      """,
      [user_id, display_name, now]
    )

    :ok
  end

  @doc "テスト用ギルドを直接DBに挿入する"
  def insert_guild(id \\ "guild_test_#{System.unique_integer([:positive])}") do
    now = DateTime.utc_now() |> DateTime.truncate(:second)

    Repo.query!(
      """
      INSERT INTO guilds (id, slug, name, description, icon, color, sort_order, created_at, updated_at)
      VALUES ($1, $1, $1, 'test guild', '🔥', '#FF0000', 0, $2, $2)
      """,
      [id, now]
    )

    id
  end

  @doc "ユーザーをギルドに参加させる（guild_memberships）"
  def insert_membership(user_id, guild_id) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)
    id = "membership_#{System.unique_integer([:positive])}"

    Repo.query!(
      """
      INSERT INTO guild_memberships (id, user_id, guild_id, joined_at, created_at, updated_at)
      VALUES ($1, $2, $3, $4, $4, $4)
      """,
      [id, user_id, guild_id, now]
    )

    :ok
  end

  @doc "チャットトークンをDBに挿入する。expires_at を省略すると1分後に設定される"
  def insert_chat_token(token_hash, user_id, guild_id, opts \\ []) do
    now = DateTime.utc_now() |> DateTime.truncate(:second)
    expires_at = Keyword.get(opts, :expires_at, DateTime.add(now, 60, :second))
    used_at = Keyword.get(opts, :used_at, nil)

    Repo.query!(
      """
      INSERT INTO chat_tokens (token_hash, user_id, guild_id, expires_at, used_at, created_at)
      VALUES ($1, $2, $3, $4, $5, $6)
      """,
      [token_hash, user_id, guild_id, expires_at, used_at, now]
    )

    :ok
  end

  @doc "SHA256 ハッシュを計算する（UserSocket と同じロジック）"
  def hash(token) do
    :crypto.hash(:sha256, token) |> Base.encode16(case: :lower)
  end
end
