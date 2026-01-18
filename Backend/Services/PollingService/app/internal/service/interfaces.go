package service

import (
	"github.com/raphael-guer1n/AREA/PollingService/internal/config"
	"github.com/raphael-guer1n/AREA/PollingService/internal/utils"
)

type ProviderConfigServiceInterface interface {
	GetProviderConfig(name string) (*config.PollingProviderConfig, error)
}

type RequestServiceInterface interface {
	ExecuteRequest(request config.PollingProviderRequestConfig, provider string, userID int, ctx utils.TemplateContext, queryOverrides map[string]string) ([]byte, error)
}
