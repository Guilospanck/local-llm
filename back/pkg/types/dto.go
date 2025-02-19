package types

type QueryRequestData struct {
	Query string `json:"query"`
}

type QueryProperty struct {
	Color    string   `json:"color"`
	SizeMin  float32  `json:"sizeMin"`
	SizeMax  float32  `json:"sizeMax"`
	PriceMin float64  `json:"priceMin"`
	PriceMax float64  `json:"priceMax"`
	Views    []string `json:"views"`
}
