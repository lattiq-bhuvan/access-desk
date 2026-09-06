package config

import (
      "github.com/lattiq/foundry/database/sql"
      "github.com/lattiq/foundry/o11y/logging"
      "github.com/lattiq/foundry/o11y/tracing"
      _http "github.com/lattiq/foundry/service/http"
)

type Config struct {
      Logging  logging.Config  `json:"logging"`
      Tracing  tracing.Config  `json:"tracing"`
      Server   _http.Config    `json:"server"`
      Database sql.Config      `json:"database"`
}