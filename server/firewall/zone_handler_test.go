package firewall

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nickrobison/terraform-linux-provider/common"
)

// mockFirewallClient is a test mock for FirewallClient interface
type mockFirewallClient struct {
	zones       []*ZoneObject
	addError    error
	getError    error
	listError   error
	deleteError error
}

func (m *mockFirewallClient) ListZones(ctx context.Context) ([]*ZoneObject, error) {
	if m.listError != nil {
		return nil, m.listError
	}
	return m.zones, nil
}

func (m *mockFirewallClient) GetZone(ctx context.Context, name string) (*ZoneObject, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	for _, z := range m.zones {
		if zName, _ := z.Name(); zName == name {
			return z, nil
		}
	}
	return nil, nil
}

func (m *mockFirewallClient) AddZone(ctx context.Context, name string, settings ZoneSettings) error {
	if m.addError != nil {
		return m.addError
	}
	return nil
}

func (m *mockFirewallClient) RemoveZone(ctx context.Context, name string) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	return nil
}

func (m *mockFirewallClient) Version() (string, error) {
	return "1.0.0", nil
}

// Tests

func TestFirewallClientInterface(t *testing.T) {
	// Test that our mock client implements the FirewallClient interface
	t.Run("Mock client implements interface", func(t *testing.T) {
		var client FirewallClient
		client = &mockFirewallClient{}

		_, err := client.Version()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		_, err = client.ListZones(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestZoneSettings(t *testing.T) {
	t.Run("ZoneSettings creation", func(t *testing.T) {
		var settings ZoneSettings
		settings.Description = "test"
		settings.Target = "default"

		if settings.Description != "test" {
			t.Errorf("Expected description 'test', got '%s'", settings.Description)
		}
	})
}

func TestHandleZoneCreate(t *testing.T) {
	client := &mockFirewallClient{}
	handler := HandleZoneCreate(client)

	reqBody := common.FirewallZoneCreateRequest{
		Name:        "testzone",
		Description: "Test Zone",
		Target:      "default",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/firewall/zones", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var resp common.FirewallZoneResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Name != "testzone" {
		t.Errorf("Expected zone name 'testzone', got '%s'", resp.Name)
	}
}

func TestHandleZoneDelete(t *testing.T) {
	client := &mockFirewallClient{}

	req := httptest.NewRequest(http.MethodDelete, "/firewall/zones/testzone", nil)
	req.SetPathValue("name", "testzone")
	w := httptest.NewRecorder()

	handler := HandleZoneDelete(client)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}
