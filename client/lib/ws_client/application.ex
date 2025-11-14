defmodule WsClient.Application do
  use Application

  def start(_type, _args) do
    children = [
      {WsClient.Supervisor, []}
    ]

    result =
      Supervisor.start_link(children,
        strategy: :one_for_one,
        name: WsClient.RootSupervisor
      )

    case result do
      {:ok, pid} ->
        maybe_start_default_client()
        {:ok, pid}

      other ->
        other
    end
  end

  defp maybe_start_default_client do
    if auto_start?() do
      default_client_opts()
      |> Keyword.put_new(:name, :ws_client_default)
      |> Keyword.put_new(:id, :ws_client_default)
      |> WsClient.start_link()
    else
      :ok
    end
  end

  defp default_client_opts do
    [
      url: endpoint(),
      interval: interval(),
      reconnect_delay: reconnect_delay()
    ]
  end

  defp endpoint do
    System.get_env("WS_CLIENT_URL", "ws://localhost:8080/ws")
  end

  defp interval do
    Application.get_env(:ws_client, :interval, 2000)
  end

  defp reconnect_delay do
    Application.get_env(:ws_client, :reconnect_delay, 1000)
  end

  defp auto_start? do
    Application.get_env(:ws_client, :auto_start, nil)
    |> case do
      nil -> Application.get_env(:ws_client, :start_connection, true)
      value -> value
    end
  end
end
