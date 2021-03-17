package iplocation

import (
	"net/http"
	"x/rest"
	"x/web"
)

const (
	URLIPGeo = "https://api.ipgeolocation.io/ipgeo?fields=geo&apiKey="
	KeyIPGeo = "29430af8568b4f609ac5fd4e00f911b5"
)

type IPGeo struct {
	IP           string `json:"ip"`
	CountryCode2 string `json:"country_code2"`
	CountryCode3 string `json:"country_code3"`
	CountryName  string `json:"country_name"`
	StateProv    string `json:"state_prov"`
	District     string `json:"district"`
	City         string `json:"city"`
	Zipcode      string `json:"zipcode"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
}
type errorGeo struct {
	Message string `json:"massage"`
}

func GetIPGeo(ip string) (*IPGeo, error) {
	var res *IPGeo
	var errGeo *errorGeo
	var url = URLIPGeo + KeyIPGeo + "ip=" + ip
	code, err := web.MethodGetNew(url, &res, &errGeo)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, rest.ErrorOK(errGeo.Message)
	}
	return res, nil
}
