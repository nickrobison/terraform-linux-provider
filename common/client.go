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
