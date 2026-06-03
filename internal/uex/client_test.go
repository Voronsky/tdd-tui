package uex_test

import (
	"testing"

	"github.com/tdd-tui/internal/uex"
)

func TestCommodityPricesAll(t *testing.T) {
	dummyClient := uex.ClientConfig("https://api.uexcorp.space/2.0", "REDACTED")
	resp, err := dummyClient.CommmodityPricesAll()
	if err != nil {
		t.Fatalf("Expected no error during unmarshal, got: %v", err)
	}

	t.Log("Data retrieved")

	if resp.HttpCode != 200 {
		t.Errorf("Expected HTTP 200, got %d", resp.HttpCode)
	}

	t.Log("Test successful HTTP Code working")

	t.Logf("Data returned: %v", resp.Data)
}
