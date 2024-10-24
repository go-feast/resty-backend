package geo

type JSONLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (d Location) ToJSON() JSONLocation {
	return JSONLocation{
		Latitude:  d.latitude,
		Longitude: d.longitude,
	}
}
