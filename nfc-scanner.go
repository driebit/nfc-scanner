package main

import (
	"bytes"
	"os"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"encoding/hex"
	"github.com/fuzxxl/nfc/2.0/nfc"
	"time"
)

func main() {
	device, err := nfc.Open("")
	if (err != nil) {
		panic(err)
	}

	access_token := get_access_token()

	var previous string
	for {
		err = device.InitiatorInit()
		if (err != nil) {
			panic(err)
		}
		modulation := nfc.Modulation{nfc.ISO14443a, nfc.Nbr106}

		t, err := device.InitiatorSelectPassiveTarget(modulation, nil)
		if (err == nfc.Error(nfc.ETIMEOUT)) {
			// On timeout, continue the loop
			continue
		} else if (err != nil) {
			panic(err)
		}

		mtarget, _ := t.(*nfc.ISO14443aTarget)
		hexUID := hex.EncodeToString(mtarget.UID[:mtarget.UIDLen])

		if (hexUID != previous) {
			fmt.Println("Read " + hexUID)
	    		register_scan(hexUID, access_token)
			fmt.Println("SENT")
		} else {
			fmt.Println("THE SAME")
			time.Sleep(100 * time.Millisecond)
		}

		previous = hexUID
	}

	nfc.Device.Close(device)
}

func logMsg(m interface{}) {
	fmt.Println(m)
}


func get_access_token() string {
	Url := os.Getenv("API_URL") + "/oauth2/token"
	Id := os.Getenv("CLIENT_ID")
	Secret := os.Getenv("CLIENT_SECRET")
	resp, err := http.PostForm(Url, url.Values{"grant_type": {"client_credentials"}, "client_id": {Id}, "client_secret": {Secret}})

	if (err != nil) {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var token_response map[string]interface{}
	err = json.Unmarshal(body, &token_response)
	if (err != nil) {
		log.Fatal(err)
	}

	return token_response["access_token"].(string)
}

func register_scan(rfid string, access_token string) {

		return

	panel_id := os.Getenv("PANEL_ID")
	type ApiRequestBody struct {
		Rfids []string `json:"rfids"`
		ObjectId string `json:"object_id"`
	}
	request_body := ApiRequestBody{
		Rfids: []string{rfid},
		ObjectId: panel_id,
	}
	json, err := json.Marshal(request_body)
	if (err != nil) {
		panic(err)
	}

	req, _ := http.NewRequest("POST", os.Getenv("API_URL") + "/tagger/actions", bytes.NewReader(json))
	req.Header.Add("Authorization", "Bearer " + access_token)
	req.Header.Add("Content-Type", "application/json2")

	client := &http.Client{}
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if (resp.StatusCode != http.StatusOK) {
		// Ignore any of the RFIDs that are not registered with Tagger
		logMsg(string(body))
	}
}
