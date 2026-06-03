package uex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Listing struct {
	ID              int     `json:"id"`
	IDCommodity     int     `json:"id_commodity"`
	IDTerminal      int     `json:"id_terminal"`
	PriceBuy        float32 `json:"price_buy"`
	PriceBuyAvg     float32 `json:"price_buy_avg"`
	PriceSell       float32 `json:"price_sell"`
	PriceSellAvg    float32 `json:"price_sell_avg"`
	SCUBuy          int     `json:"scu_buy"`
	SCUBuyAvg       int     `json:"scu_buy_avg"`
	SCUSellStock    int     `json:"scu_sell_stock"`
	SCUSellStockAvg int     `json:"scu_sell_stock_avg"`
	SCUSell         int     `json:"scu_sell"`
	SCUSellAvg      int     `json:"scu_sell_avg"`
	StatusBuy       int     `json:"status_buy,omitempty"`
	StatusSell      int     `json:"status_sell,omitempty"`
	ContainerSizes  string  `json:"container_sizes,omitempty"`
	Quality         int     `json:"quality,omitempty"`
	DateAdded       int64   `json:"date_added,omitempty"`
	DateModified    int64   `json:"date_modified,omitempty"`
	CommodityName   string  `json:"commodity_name,omitempty"`
	TerminalName    string  `json:"terminal_name,omitempty"`
}

type Commodity struct {
	IDCommodity           int    `json:"id_commodity"`
	IDStarSystem          int    `json:"id_star_system"`
	IDPlanet              int    `json:"id_planet"`
	IDOrbit               int    `json:"id_orbit"`
	IDMoon                int    `json:"id_moon"`
	IDCity                int    `json:"id_city"`
	IDOutpost             int    `json:"id_outpost"`
	IDPoi                 int    `json:"id_poi"`
	IDFaction             int    `json:"id_faction"`
	IDTerminal            int    `json:"id_terminal"`
	GameVersion           string `json:"game_version"`
	DateAdded             int    `json:"date_added"`
	DateModified          int    `json:"date_modified"`
	CommodityName         string `json:"commodity_name"`
	CommodityCode         string `json:"commodity_code"`
	CommoditySlug         string `json:"commodity_slug"`
	StarSystemName        string `json:"star_system_name"`
	PlanetName            string `json:"planet_name"`
	OrbitName             string `json:"orbit_name"`
	MoonName              string `json:"moon_name"`
	SpaceStationName      string `json:"space_station_name"`
	OutpostName           string `json:"outpost_name"`
	CityName              string `json:"city_name"`
	TerminalName          string `json:"terminal_name"`
	TerminalCode          string `json:"terminal_code"`
	TerminalSlug          string `json:"terminal_slug"`
	TerminalIsPlayerOwned int    `json:"terminal_is_player_owned"`
}

type APIResponse struct {
	Status   string    `json:"status"`
	HttpCode int       `json:"http_code"`
	Data     []Listing `json:"data"`
}

type APIClient struct {
	BaseURL string
	Token   string
	client  *http.Client
}

var (
	UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
)

func ClientConfig(url string, token string) APIClient {
	return APIClient{BaseURL: url, Token: "Bearer " + token, client: &http.Client{}}

}

func (a *APIClient) CommmodityPricesAll() (APIResponse, error) {

	req, err := http.NewRequest("GET", a.BaseURL+"/commodities_prices_all", nil)
	if err != nil {
		return APIResponse{}, fmt.Errorf("GET error = %w", err)
	}

	req.Header.Set("Authorization", a.Token)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Add("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return APIResponse{}, fmt.Errorf("Client error = %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var uexResponse APIResponse
	err = json.Unmarshal(body, &uexResponse)
	if err != nil {
		return APIResponse{}, fmt.Errorf("Client error = %w", err)
	}

	return uexResponse, nil

}

func (a *APIClient) CommmodityPrices(commodity_id int) (APIResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf(a.BaseURL+"/commodities?id_commodity=%d", commodity_id), nil)
	if err != nil {
		return APIResponse{}, fmt.Errorf("GET error = %w", err)
	}

	req.Header.Set("Authorization", a.Token)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Add("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return APIResponse{}, fmt.Errorf("Client error = %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var uexResponse APIResponse
	err = json.Unmarshal(body, &uexResponse)
	if err != nil {
		return APIResponse{}, fmt.Errorf("Client error = %w", err)
	}

	return uexResponse, nil

}

func (a *APIClient) CommmodityRoutes(src int, dest int) (APIResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf(a.BaseURL+"/commodities_route?id_terminal_origin=%d&id_terminal_destination=%d", src, dest), nil)
	if err != nil {
		return APIResponse{}, fmt.Errorf("GET error = %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.Token)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Add("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return APIResponse{}, fmt.Errorf("Client error = %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var uexResponse APIResponse
	err = json.Unmarshal(body, &uexResponse)
	if err != nil {
		return APIResponse{}, fmt.Errorf("Client error = %w", err)
	}

	return uexResponse, nil

}
