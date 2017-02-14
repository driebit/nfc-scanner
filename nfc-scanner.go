package main

import (
	"bytes"
	"os"
	"fmt"
	"log/syslog"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"encoding/hex"
	"github.com/fuzxxl/nfc/2.0/nfc"
	"time"
)

var l *syslog.Writer

func init() {
	l, _ = syslog.New(syslog.LOG_ERR, "nfc-scanner")
	defer l.Close()
}

func main() {
	access_token := get_access_token()

	// Open first available device
	device, err := nfc.Open("")
	if (err != nil) {
		panic(err)
	}

	// Discover NFC device's capabilities
	modulations, err := device.SupportedModulations(nfc.TargetMode)
	if (err != nil) {
		panic(err)
	}
	rates, err := device.SupportedBaudRates(modulations[0])
	if (err != nil) {
		panic(err)
	}
	modulation := nfc.Modulation{modulations[0], rates[0]}

	fmt.Println("Listening for NFC tags...")
	err = device.InitiatorInit()
	if (err != nil) {
		panic(err)
	}

	var previous string
	for {
		time.Sleep(100 * time.Millisecond)
		t, err := device.InitiatorSelectPassiveTarget(modulation, nil)
		if (err == nfc.Error(nfc.ETIMEOUT)) {
			// On read timeout error, just continue the loop
			continue
		} else if (err != nil) {
			// Panic on any other error
			panic(err)
		}

		if (t.String() != previous) {
			target, _ := t.(*nfc.ISO14443aTarget)
			hexUID := hex.EncodeToString(target.UID[:target.UIDLen])
			fmt.Println("Read " + hexUID)
	    		register_scan(hexUID, access_token)
		} else {
			// Delay continuation a bit to prevent hanging
			time.Sleep(100 * time.Millisecond)
		}

		previous = t.String()
	}

	nfc.Device.Close(device)
}

func get_access_token() string {
	Url := os.Getenv("API_URL") + "/oauth2/token"
	Id := os.Getenv("CLIENT_ID")
	Secret := os.Getenv("CLIENT_SECRET")
	resp, err := http.PostForm(Url, url.Values{"grant_type": {"client_credentials"}, "client_id": {Id}, "client_secret": {Secret}})

	if (err != nil) {
		panic(err)
	}

	if (resp.StatusCode != http.StatusOK) {
		panic("Could not get API access token: " + resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var token_response map[string]interface{}
	err = json.Unmarshal(body, &token_response)
	if (err != nil) {
		panic(err)
	}

	return token_response["access_token"].(string)
}

func register_scan(rfid string, access_token string) {
	object_id := os.Getenv("OBJECT_ID")
	type ApiRequestBody struct {
		Rfids []string `json:"rfids"`
		ObjectId string `json:"object_id"`
	}
	request_body := ApiRequestBody{
		Rfids: []string{rfid},
		ObjectId: object_id,
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

	if (resp.StatusCode < 200 || resp.StatusCode >= 400) {
		l.Err(resp.Status)
		// Ignore any of the RFIDs that are not registered with Tagger
		l.Err("Unregistered RFID: " + rfid + " with response: " + string(body))
	}
}
