package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type GatewayConfig struct {
	Port             int
	InternalSecret   string
	JwtPublicKey     string
	JwtPrivateKey    string
	RequestTimeoutMs int
	LogLevel         string
	DebugMode        bool
	JwtAlgorithm     string
	JwtSecret        string
	AllowedOrigins   []string
}

func LoadGatewayConfig() (*GatewayConfig, error) {
	cfg := &GatewayConfig{}

	portStr := os.Getenv("GATEWAY_PORT")
	if portStr == "" {
		return nil, fmt.Errorf("missing required env GATEWAY_PORT")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid GATEWAY_PORT: %w", err)
	}
	cfg.Port = port

	cfg.InternalSecret = os.Getenv("INTERNAL_SECRET")
	if cfg.InternalSecret == "" {
		return nil, fmt.Errorf("missing required env INTERNAL_SECRET")
	}

	cfg.JwtPublicKey = os.Getenv("JWT_PUBLIC_KEY")
	if cfg.JwtPublicKey == "" {
		return nil, fmt.Errorf("missing required env JWT_PUBLIC_KEY")
	}

	cfg.JwtPrivateKey = os.Getenv("JWT_PRIVATE_KEY")

	timeoutStr := os.Getenv("REQUEST_TIMEOUT_MS")
	if timeoutStr == "" {
		cfg.RequestTimeoutMs = 5000
	} else {
		t, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REQUEST_TIMEOUT_MS: %w", err)
		}
		cfg.RequestTimeoutMs = t
	}

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		cfg.LogLevel = lvl
	} else {
		cfg.LogLevel = "info"
	}

	debug := os.Getenv("DEBUG_MODE")
	cfg.DebugMode = (debug == "1" || debug == "true")

	cfg.JwtAlgorithm = os.Getenv("JWT_ALGO")
	if cfg.JwtAlgorithm == "" {
		cfg.JwtAlgorithm = "RS256"
	}
	cfg.JwtSecret = os.Getenv("JWT_SECRET")
	if cfg.JwtAlgorithm != "RS256" && cfg.JwtAlgorithm != "HS256" {
	    return nil, fmt.Errorf("invalid JWT_ALGO: must be RS256 or HS256")
	}

	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins != "" {
		for _, o := range strings.Split(origins, ",") {
			if trimmed := strings.TrimSpace(o); trimmed != "" {
				cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
			}
		}
	}

	return cfg, nil
}
