defmodule WsClient.Connection do
  use WebSockex

  def start_link(opts) do
    url = Keyword.fetch!(opts, :url)
    interval = Keyword.get(opts, :interval, 2000)
    state = %{interval: interval}
    WebSockex.start_link(url, __MODULE__, state, name: __MODULE__)
  end

  def handle_connect(_conn, state) do
    schedule_tick(state.interval)
    {:ok, state}
  end

  def handle_frame({:text, message}, state) do
    IO.puts("received #{message}")
    {:ok, state}
  end

  def handle_info(:tick, state) do
    schedule_tick(state.interval)
    payload = "client #{System.system_time(:millisecond)}"
    {:reply, {:text, payload}, state}
  end

  def handle_disconnect(_conn, state) do
    Process.sleep(1000)
    {:reconnect, state}
  end

  defp schedule_tick(interval) do
    Process.send_after(self(), :tick, interval)
  end
end
