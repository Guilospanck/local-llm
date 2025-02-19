package domain

type Property struct {
	Id      int     `json:"-" db:"id"`
	Color   string  `json:"color" db:"color"`
	Price   float64 `json:"price" db:"price"`
	SizeSqm float32 `json:"sizeSqm" db:"size_sqm"`
}
