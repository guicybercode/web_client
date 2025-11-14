defmodule WsClient.Application do
  use Application

  def start(_type, _args) do
    children =
      if start_connection?() do
        [{WsClient.Connection, [url: endpoint(), interval: interval()]}]
      else
        []
      end

    Supervisor.start_link(children, strategy: :one_for_one, name: WsClient.Supervisor)
  end

  defp endpoint do
    System.get_env("WS_CLIENT_URL", "ws://localhost:8080/ws")
  end

  defp interval do
    Application.get_env(:ws_client, :interval, 2000)
  end

  defp start_connection? do
    Application.get_env(:ws_client, :start_connection, true)
  end
end
