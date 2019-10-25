package web

import (
	"strconv"
	"strings"
	"x/rest"
)

type IGetable interface {
	Get(key string) string
}

func ParseFloat64(key string, g IGetable) (float64, error) {
	var v, err = strconv.ParseFloat(g.Get(key), 64)
	if err != nil {
		return 0, rest.BadRequest(key + " must be float64")
	}
	return v, nil
}

func MustGetInt64(key string, g IGetable) int64 {
	var v, err = strconv.ParseInt(g.Get(key), 10, 64)
	if err != nil {
		panic(rest.BadRequest(key + " must be int"))
	}
	return v
}

func GetArrString(key string, sep string, g IGetable) []string {
	var value = g.Get(key)
	if len(value) < 1 {
		return []string{}
	}
	return strings.Split(value, sep)
}

func GetArrInt(key string, sep string, g IGetable) ([]int, error) {
	var value = g.Get(key)
	if len(value) < 1 {
		return []int{}, rest.BadRequest(key + " must be int")
	}
	var resArr = []int{}
	var valArr = strings.Split(value, sep)
	for _, val := range valArr {
		valRes, err := strconv.Atoi(val)
		if err != nil {
			return resArr, rest.BadRequest(key + " must be int")
		}
		resArr = append(resArr, valRes)
	}
	return resArr, nil
}
