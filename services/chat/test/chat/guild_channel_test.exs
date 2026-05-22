defmodule Chat.GuildChannelTest do
  use Chat.ChannelCase

  alias Chat.UserSocket

  defp connect_as(user_id, guild_id) do
    raw_token = "token-#{System.unique_integer()}"
    insert_chat_token(hash(raw_token), user_id, guild_id)
    {:ok, socket} = connect(UserSocket, %{"token" => raw_token})
    socket
  end

  describe "join/3" do
    test "メンバーが自分のギルドに参加できる" do
      user_id = insert_user()
      guild_id = insert_guild()
      insert_membership(user_id, guild_id)
      socket = connect_as(user_id, guild_id)

      assert {:ok, %{messages: messages}, _channel_socket} =
               subscribe_and_join(socket, "guild:#{guild_id}", %{})

      assert is_list(messages)
    end

    test "join 成功時に過去メッセージが返される" do
      user_id = insert_user()
      guild_id = insert_guild()
      insert_membership(user_id, guild_id)
      insert_github_account(user_id, "github_octocat")
      insert_user_profile(user_id, "Octo Mage")

      {:ok, _} = Chat.Messages.create(guild_id, user_id, "hello")
      {:ok, _} = Chat.Messages.create(guild_id, user_id, "world")

      socket = connect_as(user_id, guild_id)

      assert {:ok, %{messages: messages}, _} =
               subscribe_and_join(socket, "guild:#{guild_id}", %{})

      bodies = Enum.map(messages, & &1.body)
      assert "hello" in bodies
      assert "world" in bodies
      assert Enum.all?(messages, &(&1.user_name == "Octo Mage"))
    end

    test "プロフィールがない場合は GitHub username が過去メッセージに返される" do
      user_id = insert_user()
      guild_id = insert_guild()
      insert_membership(user_id, guild_id)
      insert_github_account(user_id, "github_fallback")

      {:ok, _} = Chat.Messages.create(guild_id, user_id, "hello")

      socket = connect_as(user_id, guild_id)

      assert {:ok, %{messages: [message]}, _} =
               subscribe_and_join(socket, "guild:#{guild_id}", %{})

      assert message.user_name == "github_fallback"
    end

    test "別ギルドのチャンネルには参加できない" do
      user_id = insert_user()
      guild_id = insert_guild()
      other_guild_id = insert_guild()
      insert_membership(user_id, guild_id)
      socket = connect_as(user_id, guild_id)

      assert {:error, %{reason: "unauthorized"}} =
               subscribe_and_join(socket, "guild:#{other_guild_id}", %{})
    end

    test "メンバーでないユーザーは参加できない" do
      user_id = insert_user()
      guild_id = insert_guild()
      # メンバーシップを insert しない
      socket = connect_as(user_id, guild_id)

      assert {:error, %{reason: "unauthorized"}} =
               subscribe_and_join(socket, "guild:#{guild_id}", %{})
    end
  end

  describe "handle_in new_message" do
    test "メッセージを送信するとブロードキャストされる" do
      user_id = insert_user()
      guild_id = insert_guild()
      insert_membership(user_id, guild_id)
      insert_github_account(user_id, "github_octocat")
      insert_user_profile(user_id, "Octo Mage")
      socket = connect_as(user_id, guild_id)

      {:ok, _, channel_socket} =
        subscribe_and_join(socket, "guild:#{guild_id}", %{})

      ref = push(channel_socket, "new_message", %{"body" => "hello chat"})
      assert_reply(ref, :ok)

      assert_broadcast("new_message", %{
        body: "hello chat",
        user_id: ^user_id,
        user_name: "Octo Mage"
      })
    end

    test "空ボディは保存されない（バリデーションエラー）" do
      user_id = insert_user()
      guild_id = insert_guild()
      insert_membership(user_id, guild_id)
      socket = connect_as(user_id, guild_id)

      {:ok, _, channel_socket} =
        subscribe_and_join(socket, "guild:#{guild_id}", %{})

      ref = push(channel_socket, "new_message", %{"body" => ""})
      assert_reply(ref, :error, %{reason: "failed"})
    end

    test "1000文字超えは保存されない（バリデーションエラー）" do
      user_id = insert_user()
      guild_id = insert_guild()
      insert_membership(user_id, guild_id)
      socket = connect_as(user_id, guild_id)

      {:ok, _, channel_socket} =
        subscribe_and_join(socket, "guild:#{guild_id}", %{})

      long_body = String.duplicate("a", 1001)
      ref = push(channel_socket, "new_message", %{"body" => long_body})
      assert_reply(ref, :error, %{reason: "failed"})
    end
  end
end
