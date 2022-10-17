package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type DataJSON struct {
	Status `json:"status"`
}

type Level struct {
	Water string
	Wind  string
}

type Response struct {
	Status
	Level
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env")
	}
	PORT := os.Getenv("PORT")

	rand.Seed(time.Now().UnixNano())

	go randomize()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, _ := template.ParseFiles("index.html")

		b, err := ioutil.ReadFile("data.json")
		if err != nil {
			fmt.Fprint(w, "read file error")
			return
		}

		var data = DataJSON{Status: Status{}}
		if err = json.Unmarshal(b, &data); err != nil {
			fmt.Fprint(w, "marshalling error")
		}

		var response = Response{Status: data.Status}
		response.Level.Water = evaluateWater(data.Status.Water)
		response.Level.Wind = evaluateWind(data.Status.Wind)

		tpl.ExecuteTemplate(w, "index.html", response)
	})

	err = http.ListenAndServe(PORT, nil)
	log.Fatal(err)
}

func randomize() {
	for {
		var data = DataJSON{Status: Status{}}
		statusMin := 1
		statusMax := 100

		data.Status.Water = rand.Intn(statusMax-statusMin) + statusMin
		data.Status.Wind = rand.Intn(statusMax-statusMin) + statusMin

		b, err := json.MarshalIndent(&data, "", " ")
		if err != nil {
			log.Fatalln("error while marshalling json data => ", err.Error())
		}

		err = ioutil.WriteFile("data.json", b, 0644)
		if err != nil {
			log.Fatalln("error while writing value to data.json file  =>", err.Error())
		}
		time.Sleep(time.Second * 15)
	}
}

func evaluateWater(status int) string {
	if status > 8 {
		return "Bahaya"
	}

	if status > 5 {
		return "Siaga"
	}

	return "Aman"
}

func evaluateWind(status int) string {
	if status > 15 {
		return "Bahaya"
	}

	if status > 7 {
		return "Siaga"
	}

	return "Aman"
}
