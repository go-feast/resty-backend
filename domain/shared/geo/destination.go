package geo

import "errors"

// Location represents position coordinates.
type Location struct {
	latitude  float64
	longitude float64
}

func (d Location) Latitude() float64 {
	return d.latitude
}

func (d Location) Longitude() float64 {
	return d.latitude
}

var (
	ErrInvalidLatitude  = errors.New("invalid latitude: must be between -90 and 90")
	ErrInvalidLongitude = errors.New("invalid longitude: must be between -180 and 180")
)

func NewDestination(lat, long float64) (Location, error) {
	if !(lat >= -90 && lat <= 90) {
		return Location{}, ErrInvalidLatitude
	}

	if !(long >= -180 && long <= 180) {
		return Location{}, ErrInvalidLongitude
	}

	return Location{
		latitude:  lat,
		longitude: long,
	}, nil
}

func MustLocation(lat, long float64) Location {
	if !(lat >= -90 && lat <= 90) {
		panic(ErrInvalidLatitude)
	}

	if !(long >= -180 && long <= 180) {
		panic(ErrInvalidLongitude)
	}

	return Location{
		latitude:  lat,
		longitude: long,
	}
}
