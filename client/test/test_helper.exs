Application.put_env(:ws_client, :auto_start, false)
ExUnit.start()
{:ok, _} = Application.ensure_all_started(:ws_client)
