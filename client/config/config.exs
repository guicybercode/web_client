import Config

config :ws_client,
  start_connection: true,
  interval: 2000

import_config "#{config_env()}.exs"
