package main

import _ "github.com/joho/godotenv/autoload"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var url = "http://api.openweathermap.org/data/2.5/forecast?id=" + os.Getenv("CITY_ID") + "&appid=" + os.Getenv("API_KEY") + "&units=imperial"
var ifTTTAPIURL = "https://maker.ifttt.com/trigger/" + eventType + "/with/key/" + os.Getenv("IFTTT_API_KEY")

const eventType = "weather_alert"

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func jsonBody(temp float64) *bytes.Buffer {
	jsonReq := fmt.Sprintf(`{"value1":"%s"}`, fmt.Sprintf("%.2f", temp))
	fmt.Println(jsonReq)
	jsonStr := []byte(jsonReq)
	reqBody := bytes.NewBuffer(jsonStr)
	return reqBody
}

func sendSMS(temp float64) {
	req, _ := http.NewRequest("POST", ifTTTAPIURL, jsonBody(temp))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	perror(err)
	fmt.Println(resp.Status)
}

func getWeather(url string) []byte {
	resp, err := http.Get(url)
	perror(err)
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	perror(err2)
	return body
}

func main() {
	data := map[string]interface{}{}
	weather := string(getWeather(url))
	dec := json.NewDecoder(strings.NewReader(weather))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	temp, _ := jq.Float("list", "0", "main", "temp_min")
	sendSMS(temp)
}
