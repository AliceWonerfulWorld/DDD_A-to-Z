defmodule Chat.UserSocketTest do
  use Chat.ChannelCase

  alias Chat.UserSocket

  describe "connect/3" do
    test "有効なトークンで接続できる" do
      user_id = insert_user()
      guild_id = insert_guild()
      raw_token = "valid-token-#{System.unique_integer()}"
      insert_chat_token(hash(raw_token), user_id, guild_id)

      assert {:ok, socket} = connect(UserSocket, %{"token" => raw_token})
      assert socket.assigns.user_id == user_id
      assert socket.assigns.guild_id == guild_id
    end

    test "トークンは1回しか使えない（used_at が設定済みなら拒否）" do
      user_id = insert_user()
      guild_id = insert_guild()
      raw_token = "one-time-token-#{System.unique_integer()}"
      used_at = DateTime.utc_now() |> DateTime.truncate(:second)
      insert_chat_token(hash(raw_token), user_id, guild_id, used_at: used_at)

      assert :error = connect(UserSocket, %{"token" => raw_token})
    end

    test "有効期限切れトークンは拒否される" do
      user_id = insert_user()
      guild_id = insert_guild()
      raw_token = "expired-token-#{System.unique_integer()}"
      past = DateTime.add(DateTime.utc_now(), -120, :second) |> DateTime.truncate(:second)
      insert_chat_token(hash(raw_token), user_id, guild_id, expires_at: past)

      assert :error = connect(UserSocket, %{"token" => raw_token})
    end

    test "存在しないトークンは拒否される" do
      assert :error = connect(UserSocket, %{"token" => "nonexistent-token"})
    end

    test "token パラメータなしは拒否される" do
      assert :error = connect(UserSocket, %{})
    end

    test "接続後トークンは消費済みになる（再接続できない）" do
      user_id = insert_user()
      guild_id = insert_guild()
      raw_token = "consume-test-#{System.unique_integer()}"
      insert_chat_token(hash(raw_token), user_id, guild_id)

      assert {:ok, _socket} = connect(UserSocket, %{"token" => raw_token})
      assert :error = connect(UserSocket, %{"token" => raw_token})
    end
  end
end
