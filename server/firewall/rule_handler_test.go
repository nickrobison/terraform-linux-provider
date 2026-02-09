package firewall

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nickrobison/terraform-linux-provider/common"
)

func TestHandleRuleAdd(t *testing.T) {
	client := &mockFirewallClient{}
	handler := HandleRuleAdd(client)

	t.Run("Add rich rule", func(t *testing.T) {
		reqBody := common.FirewallRuleRequest{
			Zone:     "public",
			RuleType: "rich",
			Rule:     "rule family=ipv4 source address=192.168.1.0/24 accept",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/firewall/rules", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var resp common.FirewallRuleResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if resp.Zone != "public" {
			t.Errorf("Expected zone 'public', got '%s'", resp.Zone)
		}
		if resp.RuleType != "rich" {
			t.Errorf("Expected rule_type 'rich', got '%s'", resp.RuleType)
		}
	})

	t.Run("Add port rule", func(t *testing.T) {
		reqBody := common.FirewallRuleRequest{
			Zone:     "public",
			RuleType: "port",
			Port:     "8080",
			Protocol: "tcp",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/firewall/rules", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("Add service rule", func(t *testing.T) {
		reqBody := common.FirewallRuleRequest{
			Zone:     "public",
			RuleType: "service",
			Service:  "http",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/firewall/rules", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})
}

func TestHandleRuleRemove(t *testing.T) {
	client := &mockFirewallClient{}

	t.Run("Remove rich rule", func(t *testing.T) {
		reqBody := common.FirewallRuleRequest{
			Zone:     "public",
			RuleType: "rich",
			Rule:     "rule family=ipv4 source address=192.168.1.0/24 accept",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodDelete, "/firewall/rules", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler := HandleRuleRemove(client)
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
