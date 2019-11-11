package lisence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hyperboloide/lk"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func writeLicenseFile(licenseStr string) {
	// Write license file
	err_license := ioutil.WriteFile("license.dat", []byte(licenseStr), 0644)
	check(err_license)

}

func (doc *License) GenerateLicense() error {
	var start, _ = time.Parse(time.RFC3339, doc.Start.(string))
	doc.Start = start
	var end, _ = time.Parse(time.RFC3339, doc.End.(string))
	doc.End = end
	fmt.Println(start.Unix(), end)
	// Unmarshal the private key:
	privateKey, err := lk.PrivateKeyFromB64String(privateKeyStr)
	if err != nil {
		return err
	}

	// marshall the document to json bytes:
	docBytes, err := json.Marshal(&doc)
	if err != nil {
		return err

	}

	// generate your license with the private key and the document:
	license, err := lk.NewLicense(privateKey, docBytes)
	if err != nil {
		return err
	}

	// encode the new license to b64, this is what you give to your customer.
	str64, err := license.ToB64String()
	if err != nil {
		return err
	}
	fmt.Println(str64)

	writeLicenseFile(str64)
	return nil
}

func initialize() {
	// create a new Private key:
	privateKey, err := lk.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyStr, err := privateKey.ToB64String()
	fmt.Println(privateKeyStr)

	// create a license document:
	doc := License{
		// "",
		// LICENSE_NEW,
		// "VCB",
		// time.Now().Add(time.Hour * 24 * 365), // 1 year,
		// time.Now(),
		// 10,
		// 1,
		// 1,
	}

	// marshall the document to json bytes:
	docBytes, err := json.Marshal(doc)
	if err != nil {
		log.Fatal(err)

	}

	// generate your license with the private key and the document:
	license, err := lk.NewLicense(privateKey, docBytes)
	if err != nil {
		log.Fatal(err)

	}

	// encode the new license to b64, this is what you give to your customer.
	str64, err := license.ToB64String()
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println(str64)

	// get the public key. The public key should be hardcoded in your app
	// to check licences. Do not distribute the private key!
	publicKey := privateKey.GetPublicKey()
	fmt.Println("Public Key base64")
	publicKeyStr := publicKey.ToB64String()
	fmt.Println(publicKeyStr)

}
