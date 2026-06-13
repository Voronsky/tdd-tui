package uex_test

import (
	"testing"

	"github.com/tdd-tui/internal/uex"
)

func TestCommodities(t *testing.T) {
	dummyClient := uex.ClientConfig("https://api.uexcorp.space/2.0", "")
	resp, err := dummyClient.Commodities()
	if err != nil {
		t.Fatalf("Expected no error during unmarshal, got: %v", err)
	}

	t.Log("Data retrieved")

	if resp.HttpCode != 200 {
		t.Errorf("Expected HTTP 200, got %d", resp.HttpCode)
	}

	t.Log("Test successful HTTP Code working")

	t.Logf("Commodities returned: %v, %d", resp.Data, len(resp.Data))
}

//TODO: Defunc
//func TestCommodityPricesAll(t *testing.T) {
//	dummyClient := uex.ClientConfig("https://api.uexcorp.space/2.0", "")
//	resp, err := dummyClient.CommmodityPricesAll()
//	if err != nil {
//		t.Fatalf("Expected no error during unmarshal, got: %v", err)
//	}
//
//	t.Log("Data retrieved")
//
//	if resp.HttpCode != 200 {
//		t.Errorf("Expected HTTP 200, got %d", resp.HttpCode)
//	}
//
//	t.Log("Test successful HTTP Code working")
//
//	t.Logf("Commodities returned: %v, %d", resp.Data, len(resp.Data))
//}

//TODO: Defunc
//func TestCommodityAverages(t *testing.T) {
//	dummyClient := uex.ClientConfig("https://api.uexcorp.space/2.0", "")
//	avgs, err := dummyClient.CommoditiesAveragesAll()
//	if err != nil {
//		t.Fatalf("Expected no error during unmarshal, got: %v", err)
//	}
//	t.Log("Data retrieved")
//	t.Logf("Commodities Averages returned: %v", avgs)
//}
