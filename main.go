package tddtui

type Commodity struct {
	Name string
	Code string
	Kind string
	Buy  float64
	Sell float64
	SCU  float64
}

func getCommodities() ([]Commodity, error) {
	commodities := []Commodity{}
	//Rest CALL
	return commodities, nil
}
