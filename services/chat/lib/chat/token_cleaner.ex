defmodule Chat.TokenCleaner do
  use GenServer

  @interval_ms :timer.hours(1)
  # 有効期限から1時間後に削除（デバッグ用にログに残す猶予）
  @retention_after_expiry "1 hour"

  def start_link(_opts), do: GenServer.start_link(__MODULE__, [], name: __MODULE__)

  @impl true
  def init(state) do
    schedule()
    {:ok, state}
  end

  @impl true
  def handle_info(:clean, state) do
    clean()
    schedule()
    {:noreply, state}
  end

  defp schedule, do: Process.send_after(self(), :clean, @interval_ms)

  defp clean do
    {:ok, result} =
      Chat.Repo.query(
        "DELETE FROM chat_tokens WHERE expires_at < NOW() - interval '#{@retention_after_expiry}'"
      )

    if result.num_rows > 0 do
      require Logger
      Logger.info("TokenCleaner: deleted #{result.num_rows} expired chat tokens")
    end
  end
end
