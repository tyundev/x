package iplocation

import (
	"fmt"
	"strconv"
)

type Location struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Country  string  `json:"country"`
	City     string  `json:"city"`
	Address  string  `json:"address"`
	District string  `json:"district"`
}

var isUseGeo bool

func GetLocation(ip string) (*Location, error) {
	//if !isUseGeo {
	var res, err = GetIPGeo(ip)
	fmt.Println(res)
	if err != nil || res == nil {
		return nil, err
	}
	var lat, _ = strconv.ParseFloat(res.Latitude, 64)
	var lng, _ = strconv.ParseFloat(res.Longitude, 64)
	var city = res.City
	if city == "" {
		city = res.StateProv
	}
	var l = &Location{
		Lat:      lat,
		Lng:      lng,
		Country:  res.District + ", " + city + ", " + res.CountryName,
		City:     city,
		District: res.District,
	}
	isUseGeo = false
	return l, err
	// }
	// var res, err = GetIPStack(ip)
	// if err != nil {
	// 	return nil, err
	// }
	// var l = &Location{
	// 	Lat:     res.Latitude,
	// 	Lng:     res.Longitude,
	// 	Country: res.CountryName,
	// 	City:    res.City,
	// }
	// isUseGeo = true
	// return l, err
}

func ResIPStack(ip string) (*Location, error) {
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
	return l, err
}
