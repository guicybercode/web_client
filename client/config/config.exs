import Config

config :ws_client,
  auto_start: true,
  interval: 2000,
  reconnect_delay: 1000

import_config "#{config_env()}.exs"
