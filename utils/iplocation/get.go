package iplocation

import "strconv"

type Location struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Country string  `json:"country"`
	City    string  `json:"city"`
	Address string  `json:"address"`
}

var isUseGeo bool

func GetLocation(ip string) (*Location, error) {
	if !isUseGeo {
		var res, err = GetIPGeo(ip)
		if err != nil {
			return nil, err
		}
		var lat, _ = strconv.ParseFloat(res.Latitude, 64)
		var lng, _ = strconv.ParseFloat(res.Longitude, 64)
		var l = &Location{
			Lat:     lat,
			Lng:     lng,
			Country: res.CountryName,
			City:    res.City,
		}
		isUseGeo = true
		return l, err
	}
	var res, err = GetIPStack(ip)
	if err != nil {
		return nil, err
	}
	var l = &Location{
		Lat:     res.Latitude,
		Lng:     res.Longitude,
		Country: res.CountryName,
		City:    res.City,
	}
	isUseGeo = true
	return l, err
}
