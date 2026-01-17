package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/raphael-guer1n/AREA/PollingService/internal/config"
	"github.com/raphael-guer1n/AREA/PollingService/internal/utils"
)

type RequestService struct {
	oauth2TokenSvc *OAuth2TokenService
	client         *http.Client
	logRequests    bool
}

func NewRequestService(oauth2TokenSvc *OAuth2TokenService, logRequests bool) *RequestService {
	return &RequestService{
		oauth2TokenSvc: oauth2TokenSvc,
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		logRequests: logRequests,
	}
}

func (s *RequestService) ExecuteRequest(request config.PollingProviderRequestConfig, provider string, userID int, ctx utils.TemplateContext, queryOverrides map[string]string) ([]byte, error) {
	urlValue, err := utils.RenderTemplateString(request.URLTemplate, ctx)
	if err != nil {
		return nil, err
	}
	urlStr, ok := urlValue.(string)
	if !ok {
		urlStr = fmt.Sprint(urlValue)
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(parsedURL.Scheme) {
	case "webcal", "webcals":
		parsedURL.Scheme = "https"
	}
	query := parsedURL.Query()
	for key, value := range request.QueryParams {
		if strings.TrimSpace(key) == "" {
			continue
		}
		renderedValue, err := renderTemplateStringNested(value, ctx)
		if err != nil {
			var missing utils.MissingTemplateValueError
			if errors.As(err, &missing) {
				continue
			}
			return nil, err
		}
		query.Set(key, fmt.Sprint(renderedValue))
	}
	for key, value := range queryOverrides {
		if strings.TrimSpace(key) == "" {
			continue
		}
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()
	urlStr = parsedURL.String()

	var body io.Reader
	contentType := ""
	if len(request.BodyTemplate) > 0 {
		var template any
		if err := json.Unmarshal(request.BodyTemplate, &template); err != nil {
			return nil, fmt.Errorf("invalid body template: %w", err)
		}
		rendered, err := utils.RenderTemplateValue(template, ctx)
		if err != nil {
			return nil, err
		}

		switch strings.ToLower(strings.TrimSpace(request.BodyEncoding)) {
		case "", "json":
			payload, err := json.Marshal(rendered)
			if err != nil {
				return nil, fmt.Errorf("marshal request payload: %w", err)
			}
			body = bytes.NewReader(payload)
			contentType = "application/json"
		case "form", "x-www-form-urlencoded":
			form, err := buildFormPayload(rendered)
			if err != nil {
				return nil, err
			}
			body = strings.NewReader(form.Encode())
			contentType = "application/x-www-form-urlencoded"
		default:
			return nil, fmt.Errorf("unsupported body_encoding %q", request.BodyEncoding)
		}
	}

	req, err := http.NewRequest(request.Method, urlStr, body)
	if err != nil {
		return nil, err
	}

	for key, value := range request.Headers {
		renderedValue, err := renderTemplateStringNested(value, ctx)
		if err != nil {
			return nil, err
		}
		req.Header.Set(key, fmt.Sprint(renderedValue))
	}
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	if request.Auth != nil {
		switch request.Auth.Type {
		case "oauth2":
			providerName := provider
			if request.Auth.Provider != "" {
				providerName = request.Auth.Provider
			}
			token, err := s.oauth2TokenSvc.GetProviderToken(userID, providerName)
			if err != nil {
				return nil, err
			}
			prefix := request.Auth.Prefix
			req.Header.Set(request.Auth.Header, prefix+token)
		default:
			return nil, fmt.Errorf("unsupported auth type %s", request.Auth.Type)
		}
	}

	start := time.Now()
	resp, err := s.client.Do(req)
	if err != nil {
		if s.logRequests {
			log.Printf("polling: provider request failed provider=%s method=%s url=%s err=%v", provider, request.Method, urlStr, err)
		}
		return nil, fmt.Errorf("provider request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	if s.logRequests {
		log.Printf("polling: provider request provider=%s method=%s url=%s status=%d duration=%s", provider, request.Method, urlStr, resp.StatusCode, time.Since(start))
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("provider request failed: status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func renderTemplateStringNested(value string, ctx utils.TemplateContext) (any, error) {
	rendered, err := utils.RenderTemplateString(value, ctx)
	if err != nil {
		return nil, err
	}
	str, ok := rendered.(string)
	if !ok {
		return rendered, nil
	}
	if !strings.Contains(str, "{{") || !strings.Contains(str, "}}") {
		return str, nil
	}
	return utils.RenderTemplateString(str, ctx)
}

func buildFormPayload(rendered any) (url.Values, error) {
	obj, ok := rendered.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("form body must be an object")
	}

	values := url.Values{}
	for key, value := range obj {
		switch v := value.(type) {
		case []any:
			for _, item := range v {
				addFormValue(values, key, item)
			}
		case []string:
			for _, item := range v {
				addFormValue(values, key, item)
			}
		default:
			addFormValue(values, key, v)
		}
	}

	return values, nil
}

func addFormValue(values url.Values, key string, value any) {
	if value == nil {
		return
	}
	str := strings.TrimSpace(fmt.Sprint(value))
	if str == "" {
		return
	}
	values.Add(key, str)
}
