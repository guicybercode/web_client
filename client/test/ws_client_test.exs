defmodule WsClientTest do
  use ExUnit.Case

  test "configures interval" do
    assert Keyword.get([interval: 1000], :interval) == 1000
  end
end
