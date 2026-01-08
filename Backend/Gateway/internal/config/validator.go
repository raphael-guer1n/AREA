package config

import (
	"errors"
	"fmt"
	"strings"
)

var validHTTPMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"PATCH":   true,
	"DELETE":  true,
	"OPTIONS": true,
	"HEAD":    true,
}

func ValidateAll(services []ServiceConfig) error {
	for _, svc := range services {
		if err := validateServiceConfig(svc); err != nil {
			return fmt.Errorf("service '%s' invalid: %w", svc.Name, err)
		}
	}
	return nil
}

func validateServiceConfig(svc ServiceConfig) error {
	if svc.Name == "" {
		return errors.New("missing service name")
	}
	if svc.BaseURL == "" {
		return errors.New("missing base_url")
	}
	if len(svc.Routes) == 0 {
		return errors.New("service has no routes")
	}

	routeSet := make(map[string]struct{})

	for _, r := range svc.Routes {
		if err := validateRouteConfig(r); err != nil {
			return fmt.Errorf("invalid route '%s': %w", r.Path, err)
		}

		for _, m := range r.Methods {
			key := fmt.Sprintf("%s#%s", r.Path, m)
			if _, exists := routeSet[key]; exists {
				return fmt.Errorf("duplicate route+method: %s %s", m, r.Path)
			}
			routeSet[key] = struct{}{}
		}
	}

	return nil
}

func validateRouteConfig(r RouteConfig) error {
	if r.Path == "" {
		return errors.New("path cannot be empty")
	}
	if !strings.HasPrefix(r.Path, "/") {
		return errors.New("path must start with '/'")
	}
	if strings.Contains(r.Path, "//") {
		return errors.New("path cannot contain double slashes '//'")
	}

	if len(r.Methods) == 0 {
		return errors.New("route must declare at least one HTTP method")
	}

	for _, m := range r.Methods {
		if !validHTTPMethods[strings.ToUpper(m)] {
			return fmt.Errorf("invalid HTTP method '%s'", m)
		}
	}

	if r.AuthRequired && len(r.Permissions) == 0 {
		// TO DO
	}

	if r.InternalOnly && r.AuthRequired {
		return errors.New("internal_only routes cannot have auth_required=true")
	}

	return nil
}
