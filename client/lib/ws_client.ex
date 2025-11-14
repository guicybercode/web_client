defmodule WsClient do
  def start_link(opts \\ []) do
    DynamicSupervisor.start_child(WsClient.Supervisor, {WsClient.Connection, opts})
  end

  def stop(pid) do
    DynamicSupervisor.terminate_child(WsClient.Supervisor, pid)
  end
end
