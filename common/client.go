package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client *http.Client
	host   string
	port   int
}

func NewClient(host string) *Client {
	return &Client{
		client: &http.Client{},
		host:   host,
		port:   8080,
	}
}

func (c *Client) WithPort(port int) *Client {
	c.port = port
	return c
}

// ZFS

func (c *Client) ZfsCreatePool(ctx context.Context, create ZpoolCreateRequest) (ZPoolResponse, error) {
	var result ZPoolResponse
	body, err := json.Marshal(create)
	if err != nil {
		return result, err
	}
	url := c.createUrl("zfs", "zpool")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return result, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		e := fmt.Errorf("failed to create user. error: %s", string(b))
		return result, e
	}
	DecodeInto(resp, &result)
	return result, nil
}

func (c *Client) ZfsGetPools(ctx context.Context) (ZpoolListResponse, error) {

	var result ZpoolListResponse
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.createUrl("zfs", "zpool"), nil)
	if err != nil {
		return result, err
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return result, err
	}

	err = DecodeInto[ZpoolListResponse](resp, &result)
	return result, err

}

func (c *Client) ZfsGetPool(ctx context.Context, name string) (ZPoolResponse, error) {
	var pool ZPoolResponse
	url := fmt.Sprintf("%s/%s", c.createUrl("zfs", "zpool"), name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return pool, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return pool, err
	}

	err = DecodeInto[ZPoolResponse](resp, &pool)
	return pool, err
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	return c.client.Do(req)
}

func (c *Client) createUrl(module string, resource string) string {
	return fmt.Sprintf("http://%s:%d/%s/%s", c.host, c.port, module, resource)
}

// Firewall

func (c *Client) FirewallCreateZone(ctx context.Context, create FirewallZoneCreateRequest) (FirewallZoneResponse, error) {
	var result FirewallZoneResponse
	body, err := json.Marshal(create)
	if err != nil {
		return result, err
	}
	url := c.createUrl("firewall", "zones")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return result, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		e := fmt.Errorf("failed to create zone. error: %s", string(b))
		return result, e
	}
	DecodeInto(resp, &result)
	return result, nil
}

func (c *Client) FirewallGetZones(ctx context.Context) ([]FirewallZoneResponse, error) {
	var result []FirewallZoneResponse
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.createUrl("firewall", "zones"), nil)
	if err != nil {
		return result, err
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return result, err
	}

	err = DecodeInto[[]FirewallZoneResponse](resp, &result)
	return result, err
}

func (c *Client) FirewallGetZone(ctx context.Context, name string) (FirewallZoneResponse, error) {
	var zone FirewallZoneResponse
	url := fmt.Sprintf("%s/%s", c.createUrl("firewall", "zones"), name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return zone, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return zone, err
	}

	err = DecodeInto[FirewallZoneResponse](resp, &zone)
	return zone, err
}

func (c *Client) FirewallDeleteZone(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/%s", c.createUrl("firewall", "zones"), name)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete zone. error: %s", string(b))
	}

	return nil
}

func (c *Client) FirewallAddRule(ctx context.Context, rule FirewallRuleRequest) (FirewallRuleResponse, error) {
	var result FirewallRuleResponse
	body, err := json.Marshal(rule)
	if err != nil {
		return result, err
	}
	url := c.createUrl("firewall", "rules")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return result, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusCreated {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		e := fmt.Errorf("failed to add rule. error: %s", string(b))
		return result, e
	}
	DecodeInto(resp, &result)
	return result, nil
}

func (c *Client) FirewallRemoveRule(ctx context.Context, rule FirewallRuleRequest) error {
	body, err := json.Marshal(rule)
	if err != nil {
		return err
	}
	url := c.createUrl("firewall", "rules")
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove rule. error: %s", string(b))
	}

	return nil
}
