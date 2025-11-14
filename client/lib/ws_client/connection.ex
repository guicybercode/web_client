defmodule WsClient.Connection do
  use WebSockex

  def start_link(opts) do
    url = Keyword.fetch!(opts, :url)
    interval = Keyword.get(opts, :interval, 2000)
    reconnect_delay = Keyword.get(opts, :reconnect_delay, 1000)
    client_id = Keyword.get(opts, :client_id, generate_client_id())
    state = %{interval: interval, reconnect_delay: reconnect_delay, client_id: client_id}
    WebSockex.start_link(url, __MODULE__, state, connection_options(opts))
  end

  def child_spec(opts) do
    id = opts |> Keyword.get(:id) |> default_child_id(opts)

    %{
      id: id,
      start: {__MODULE__, :start_link, [opts]},
      restart: :transient,
      shutdown: 5000,
      type: :worker
    }
  end

  def handle_connect(_conn, state) do
    schedule_tick(state.interval)
    {:ok, state}
  end

  def handle_frame({:text, message}, state) do
    case Jason.decode(message) do
      {:ok, %{"client_id" => sender_id, "content" => content, "timestamp" => _timestamp}} ->
        if sender_id != state.client_id do
          IO.puts("[#{sender_id}] #{content}")
        end

      {:error, _} ->
        IO.puts("received raw: #{message}")
    end
    {:ok, state}
  end

  def handle_info(:tick, state) do
    schedule_tick(state.interval)
    payload = "client #{System.system_time(:millisecond)}"
    {:reply, {:text, payload}, state}
  end

  def handle_disconnect(_status, state) do
    Process.sleep(state.reconnect_delay)
    {:reconnect, state}
  end

  defp schedule_tick(interval) do
    Process.send_after(self(), :tick, interval)
  end

  defp connection_options(opts) do
    base = [async: true, handle_initial_conn_failure: true]

    case Keyword.get(opts, :name) do
      nil -> base
      name -> Keyword.put(base, :name, name)
    end
  end

  defp default_child_id(nil, opts) do
    Keyword.get(opts, :name) || {__MODULE__, System.unique_integer([:positive, :monotonic])}
  end

  defp default_child_id(id, _opts), do: id

  defp generate_client_id do
    :crypto.strong_rand_bytes(8) |> Base.encode16(case: :lower)
  end
end
