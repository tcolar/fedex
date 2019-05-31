// main uploads letterhead.png and signature.png to FedEx prod
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/happyreturns/fedex"
	"github.com/happyreturns/fedex/models"
)

func main() {

	happyReturnsLetterhead := base64EncodedFile("happyreturns_letterhead.png")
	cutsClothingLetterhead := base64EncodedFile("cutsclothing_letterhead.png")
	signature := base64EncodedFile("signature.png")

	credData, err := ioutil.ReadFile("../creds.json")
	if err != nil {
		panic(err)
	}

	creds := map[string]fedex.Fedex{}
	if err := json.Unmarshal(credData, &creds); err != nil {
		panic(err)
	}

	prodFedex := creds["prod"]

	err = prodFedex.UploadImages([]models.Image{
		{
			ID:    "IMAGE_1",
			Image: happyReturnsLetterhead,
		},
		{
			ID:    "IMAGE_2",
			Image: signature,
		},
		{
			ID:    "IMAGE_3",
			Image: cutsClothingLetterhead,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("DONE")
}

//  ...
func base64EncodedFile(filename string) string {
	readBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(readBytes)
}
