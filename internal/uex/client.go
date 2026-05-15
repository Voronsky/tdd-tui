package uex

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Listing struct {
	ID              int     `json:"id"`
	IDCommodity     int     `json:"id_commodity"`
	IDTerminal      int     `json:"id_terminal"`
	PriceBuy        float32 `json:"price_buy"`
	PriceBuyAvg     float32 `json:"price_buy_avg"`
	SCUBuy          int     `json:"scu_buy"`
	SCUBuyAvg       int     `json:"scu_buy_avg"`
	SCUSellStock    int     `json:"scu_sell_stock"`
	SCUSellStockAvg int     `json:"scu_sell_stock_avg"`
	SCUSell         int     `json:"scu_sell"`
	SCUSellAvg      int     `json:"scu_sell_avg"`
	StatusBuy       int     `json:"status_buy"`
	StatusSell      int     `json:"status_sell"`
	ContainerSizes  string  `json:"container_sizes"`
	Quality         int     `json:"quality"`
	DateAdded       int64   `json:"date_added"`
	DateModified    int64   `json:"date_modified"`
	CommodityName   string  `json:"commodity_name"`
	TerminalName    string  `json:"terminal_name"`
}

type APIResponse struct {
	Status   string    `json:"status"`
	HttpCode int       `json:"http_code"`
	Data     []Listing `json:"data"`
}

type APIClient struct {
	BaseURL string
	Token   string
}

func ClientConfig(url string, token string) APIClient {
	return APIClient{BaseURL: url, Token: "Bearer " + token}

}

func (a *APIClient) CommmoddityPrices() (APIResponse, error) {

	req, err := http.NewRequest("GET", a.BaseURL+"/commodities", nil)
	if err != nil {
		log.Panicln("Error unable to make a New Request")
		return APIResponse{}, err
	}

	req.Header.Set("Authorization", a.Token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Get request failed lol")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var uexResponse APIResponse
	err = json.Unmarshal(body, &uexResponse)
	if err != nil {
		log.Fatalln("Unable to Marshal UEX response")
	}

	return uexResponse, nil

}

func (a *APIClient) CommmoddityRoutes(src int, dest int) (APIResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf(a.BaseURL+"/commodities_route?id_terminal_origin=%d&id_terminal_destination=%d", src, dest), nil)
	if err != nil {
		log.Panicln("Error unable to make a New Request")
		return APIResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+a.Token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Get request failed lol")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var uexResponse APIResponse
	err = json.Unmarshal(body, &uexResponse)
	if err != nil {
		log.Fatalln("Unable to Marshal UEX response")
	}

	return uexResponse, nil

}
