package crc

import (
	"fmt"
	"math/big"
	"strconv"
	"x/rest"
)

const (
	AllNumByte_Status = "00BD"
	CMD_Status        = "03"
	ADDR_Status       = "1000"
	NumDATA_Status    = "0020"

	AllNumByte_Setup = "006D"
	CMD_Setup        = "10"
	ADDR_Setup       = "6000"
	NumDATA_Setup    = "0010"
)

func buildTextToByte(data string) []byte {
	fmt.Println("DATA: ", data)
	var countData = len(data)
	var temp = make([]string, countData)
	for i, s := range data {
		temp[i] = string(s)
	}
	var dataByte = make([]byte, countData*4)
	for i := 0; i < countData; i++ {

		var t = []byte(temp[i])
		for j := 0; j < len(t); j++ {
			dataByte[j*4] = t[j]
		}
	}
	return dataByte
}

func ConvertHexStrToDec(hex string) (int, error) {
	i := new(big.Int)
	if _, isOK := i.SetString(hex, 16); isOK {
		return strconv.Atoi(i.String())
	}
	return 0, rest.ErrorOK("Lỗi convert")
}

func convertDataToFormat(value int) (float64, error) {
	var floatStr = strconv.Itoa(value)
	if len(floatStr) >= 2 {
		runes := []rune(floatStr)
		floatStr = string(runes[0:(len(runes)-2)]) + "." + string(runes[len(runes)-1])
	}
	return strconv.ParseFloat(floatStr, 64)
}

func ConvertHexStrToFormat(hex string) (float64, error) {
	var dataInt, err = ConvertHexStrToDec(hex)
	if err != nil {
		return 0, err
	}
	return convertDataToFormat(dataInt)
}

func GetCrcConvert(arrValue []string) string {
	var lenArrStop = len(arrValue) - 2
	var valueConvertCrc = ""
	for i, val := range arrValue {
		if i == lenArrStop {
			break
		}
		valueConvertCrc += val + ","
	}
	return valueConvertCrc
}

func ConvertIntToHex(num int) string {
	h := fmt.Sprintf("%X", num)
	var valueO = 4 - len(h)
	var addValueO = ""
	if valueO > 0 {
		for i := 0; i < valueO; i++ {
			addValueO += "0"
		}
	}
	h = addValueO + h
	return h
}
