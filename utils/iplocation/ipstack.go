package iplocation

import (
	"x/rest"
	"x/web"
)

const (
	URLIPStack = "http://api.ipstack.com/"
	KeyIPStack = "4ee5262f6cbd502aef84d6dbb027e25c"
)

type IPStack struct {
	IP            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      struct {
		GeonameID int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag             string `json:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code"`
		IsEu                    bool   `json:"is_eu"`
	} `json:"location"`
	TimeZone struct {
		ID               string `json:"id"`
		CurrentTime      string `json:"current_time"`
		GmtOffset        int    `json:"gmt_offset"`
		Code             string `json:"code"`
		IsDaylightSaving bool   `json:"is_daylight_saving"`
	} `json:"time_zone"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Plural       string `json:"plural"`
		Symbol       string `json:"symbol"`
		SymbolNative string `json:"symbol_native"`
	} `json:"currency"`
	Connection struct {
		Asn int    `json:"asn"`
		Isp string `json:"isp"`
	} `json:"connection"`
}

type dataResponse struct {
	*IPStack `bson:",inline"`
	Success  bool `json:"success"`
	Error    struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}

func GetIPStack(ip string) (*IPStack, error) {
	var res *dataResponse
	var url = URLIPStack + ip + "access_key=" + KeyIPStack
	_, err := web.MethodGet(url, &res)
	if err != nil {
		return nil, err
	}
	if !res.Success {
		return nil, rest.ErrorOK(res.Error.Info)
	}
	return res.IPStack, nil
}

func GetIPStacks(ips []string) ([]*IPStack, error) {
	var res []*IPStack
	var ipAll = ""
	for i, val := range ips {
		if i == 0 {
			ipAll = val
		} else {
			ipAll += "," + val
		}
	}
	var url = URLIPStack + ipAll + "access_key=" + KeyIPStack
	_, err := web.MethodGet(url, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
