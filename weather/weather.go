package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Weather struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
	Weather []Weather
	Clouds  struct {
		All int `json:"all"`
	}
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	}
	Sys struct {
		Type    int     `json:"type"`
		Id      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int64   `json:"sunrise"`
		Sunset  int64   `json:"sunset"`
	}
	Base       string `json:"base"`
	Visibility int    `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   float64 `json:"deg"`
		Gust  float64 `json:"gust"`
	}
	Dt   int64  `json:"dt"`
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

func main() {
	// flags
	var apiKey string
	var cityId string
	flag.StringVar(&apiKey, "key", "YourKeyFromOpenWeatherMapsDotCom", "API Key from Openweather")
	flag.StringVar(&cityId, "city", "1111111", "City ID fromOpenweather")
	flag.Parse()

	var weather WeatherData
	var owmUrl string = "http://api.openweathermap.org/data/2.5/weather"
	var units string = "imperial"
	var weatherUrl string = owmUrl + "?id=" + cityId + "&appid=" + apiKey + "&units=" + units

	var myClient = &http.Client{Timeout: 10 * time.Second}
	// create reader to get URL request
	response, err := myClient.Get(weatherUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	decodeErr := json.Unmarshal(body, &weather)
	if decodeErr != nil {
		log.Fatal(decodeErr)
	}
	currentTemp := strconv.FormatFloat(weather.Main.Temp, 'f', -1, 64)
	currentDesc := weather.Weather[0].Description
	retVal := currentTemp + "°F : " + currentDesc
	fmt.Print(retVal)
}
