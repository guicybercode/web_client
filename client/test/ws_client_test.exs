defmodule WsClientTest do
  use ExUnit.Case, async: false

  @unreachable "ws://localhost:65535/ws"

  test "connection child_spec honors custom id" do
    spec = WsClient.Connection.child_spec(url: @unreachable, id: :custom_id)
    assert spec.id == :custom_id
  end

  test "multiple clients can run concurrently" do
    {:ok, pid1} =
      WsClient.start_link(
        url: @unreachable,
        interval: 10,
        reconnect_delay: 10,
        name: :client_one,
        id: :client_one
      )

    {:ok, pid2} =
      WsClient.start_link(
        url: @unreachable,
        interval: 10,
        reconnect_delay: 10,
        name: :client_two,
        id: :client_two
      )

    refute pid1 == pid2
    Process.sleep(20)
    assert Process.alive?(pid1)
    assert Process.alive?(pid2)
    :ok = WsClient.stop(pid1)
    :ok = WsClient.stop(pid2)
  end
end
