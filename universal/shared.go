package universal


import (
  "TermoChat/config"
)


var secretKey = config.LoadEnv().SECRET
