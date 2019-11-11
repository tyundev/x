package lisence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func Read() *License {
	// Open our jsonFile
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result *License
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func (ls *License) Println() {
	fmt.Println("===== THONG TIN LICENSE ======")
	fmt.Println("== Machine      : ", ls.DeviceID)
	fmt.Println("== Type         : ", ls.Type)
	fmt.Println("== Organization : ", ls.Organization)
	fmt.Println("== Start        : ", ls.Start)
	fmt.Println("== End          : ", ls.End)
	// fmt.Println("== Counter Limit: ", ls.CounterLimit)
	// fmt.Println("== Screen Limit : ", ls.ScreenLimit)
	// fmt.Println("== Kiosk Limit  : ", ls.KioskLimit)
}
