package main

import _ "github.com/joho/godotenv/autoload"

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var url = "http://api.openweathermap.org/data/2.5/forecast?id=" + os.Getenv("CITY_ID") + "&appid=" + os.Getenv("API_KEY") + "&units=imperial"

func perror(err error) {
	if err != nil {
		panic(err)
	}
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
	fmt.Println(fmt.Sprintf("The temperature will be %.2f", temp))
}
