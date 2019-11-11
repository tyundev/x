package lisence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/denisbrodbeck/machineid"
	"github.com/golang/glog"
	"github.com/hyperboloide/lk"
)

func Decode() *License {
	var machineID, err = getMachineID()
	if err != nil {
		glog.Errorf("get Machine", err)
		return nil
	}
	return validLicense(machineID)
}

func getMachineID() (string, error) {
	return machineid.ID()
}

func validLicense(machineID string) *License {

	/*	defer func() {
		glog.Info("Invalid License file. Please contact Miraway for support")
		os.Exit(0)
	}()*/

	if !checkFileExist("license.dat") {
		glog.Info("License file not found. Please contact Miraway for support")
		//os.Exit(0)
	}

	// Read License file
	pwd, _ := os.Getwd()
	filepath := path.Join(pwd, "license.dat")
	glog.Info("File Path", filepath)

	buff1, err := ioutil.ReadFile(filepath)
	check(err)
	licenseStr := string(buff1)

	license, _ := lk.LicenseFromB64String(licenseStr)
	publicKey, _ := lk.PublicKeyFromB64String(publicKeyStr)
	if ok, err := license.Verify(publicKey); err != nil {
		glog.Info("Invalid license")
		return nil
	} else if ok {
		var temp *License
		json.Unmarshal(license.Data, &temp)

		switch {
		case temp.Type == LICENSE_NEW && temp.DeviceID == "":
			// first time license will be used
			temp.DeviceID = machineID
			temp.Type = LICENSE_USED

			// Unmarshal the private key:
			privateKey, err := lk.PrivateKeyFromB64String(privateKeyStr)
			if err != nil {
				glog.Fatal(err)
			}

			// marshall the document to json bytes:
			docBytes, err := json.Marshal(&temp)
			if err != nil {
				glog.Fatal(err)

			}
			// generate your license with the private key and the document:
			license, err := lk.NewLicense(privateKey, docBytes)
			if err != nil {
				glog.Fatal(err)

			}

			// encode the new license to b64, this is what you give to your customer.
			str64, err := license.ToB64String()
			if err != nil {
				glog.Info(err)
			}
			writeLicenseFile(str64)

			glog.Info("License applied")
		case temp.Type == LICENSE_USED && temp.DeviceID == machineID:
			// Valide license content
			glog.Info("Valided license")
		default:
			// invalide license content
			glog.Info("Invalid license")
			return nil
		}
		fmt.Println(temp)
		return temp
	}
	return nil
}

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
