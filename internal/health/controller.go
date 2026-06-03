package health

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type Controller func(w http.ResponseWriter, r *http.Request)

type Config struct {
	Healthcheck string
	Version     string
	DeployedAt  string
	Service     string
}

type Response struct {
	Status      string `json:"status"`
	Healthcheck string `json:"healthcheck"`
	Version     string `json:"version"`
	DeployedAt  string `json:"deployed_at"`
	Service     string `json:"service,omitempty"`
}

func Handler(cfg Config) Controller {
	resolved := resolveConfig(cfg)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Status:      "ok",
			Healthcheck: resolved.Healthcheck,
			Version:     resolved.Version,
			DeployedAt:  resolved.DeployedAt,
			Service:     resolved.Service,
		})
	}
}

func resolveConfig(cfg Config) Config {
	if strings.TrimSpace(cfg.Healthcheck) == "" {
		cfg.Healthcheck = "stable"
	}
	if strings.TrimSpace(cfg.Version) == "" {
		cfg.Version = "dev"
	}
	if strings.TrimSpace(cfg.DeployedAt) == "" {
		cfg.DeployedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if strings.TrimSpace(cfg.Service) == "" {
		cfg.Service = "api"
	}

	return cfg
}
