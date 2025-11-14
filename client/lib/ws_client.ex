defmodule WsClient do
  def start_link(opts \\ []) do
    WsClient.Connection.start_link(opts)
  end
end
