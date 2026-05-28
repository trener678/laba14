package sharding

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type EtcdCoordinator struct {
	endpoint string
	prefix   string
	client   *http.Client
}

func NewEtcdCoordinator(endpoints []string, prefix string, ttl int64) (*EtcdCoordinator, error) {
	if prefix == "" {
		prefix = "/laba14/collectors"
	}
	if len(endpoints) == 0 || endpoints[0] == "" {
		return nil, fmt.Errorf("at least one etcd endpoint is required")
	}
	endpoint := strings.TrimRight(endpoints[0], "/")
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return nil, err
	}
	_ = ttl
	return &EtcdCoordinator{
		endpoint: endpoint,
		prefix:   prefix,
		client:   &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (c *EtcdCoordinator) Close() error {
	return nil
}

func (c *EtcdCoordinator) Register(ctx context.Context, collectorID string) error {
	payload := map[string]string{
		"key":   encode(c.prefix + "/" + collectorID),
		"value": encode(collectorID),
	}
	return c.post(ctx, "/v3/kv/put", payload, nil)
}

func (c *EtcdCoordinator) Collectors(ctx context.Context) ([]string, error) {
	payload := map[string]string{
		"key":       encode(c.prefix + "/"),
		"range_end": encode(c.prefix + "0"),
	}
	var response struct {
		Kvs []struct {
			Key string `json:"key"`
		} `json:"kvs"`
	}
	if err := c.post(ctx, "/v3/kv/range", payload, &response); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(response.Kvs))
	for _, kv := range response.Kvs {
		keyBytes, err := base64.StdEncoding.DecodeString(kv.Key)
		if err != nil {
			continue
		}
		id := strings.TrimPrefix(string(keyBytes), c.prefix+"/")
		if id != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no active collectors under %s", c.prefix)
	}
	sort.Strings(ids)
	return ids, nil
}

func (c *EtcdCoordinator) post(ctx context.Context, path string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("etcd returned %s", resp.Status)
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
