package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var port = "81"

type response struct{
	Slack_name  		string
	Track 				string
	Current_day 		string
	Utc_time 			string
	Github_file_url 	string
	Github_repo_url 	string
	Status_code 		uint
}

func GetJSONData(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("slack_name")
	track	 := r.URL.Query().Get("track")
	
	if len(username) == 0 || len(track) == 0 {
		fmt.Println("Slack_name and track cannot be empty")
		return
	}

	var response response

	executable_file, err := os.Executable()

	if err != nil{
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	utc_time := time.Now().UTC().String()

	utcArr := strings.Split(utc_time, " ")
	timeSplit := strings.Split(utcArr[1], ".")
	formattedTime := utcArr[0] + "T" + timeSplit[0] + "Z"

	fmt.Println(wd)
	response.Slack_name = username
	response.Track = track
	response.Current_day = time.Now().Weekday().String()
	response.Utc_time = formattedTime
	response.Github_file_url = executable_file
	response.Github_repo_url = "https://github.com/ntekim/hng-stage-1"
	response.Status_code   = http.StatusOK

	resp, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		panic(err)
	}

}


func route() http.Handler{
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Content-Type", "X-CSRF-Token", "Accept"},
		AllowCredentials: true,
		MaxAge: 300,
		ExposedHeaders: []string{"Link"},
	}))

	mux.Use(middleware.Logger)
	mux.Get("/api", GetJSONData)

	return mux
}

func main() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: route(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}