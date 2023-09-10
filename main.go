package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var port = "81"

type responseData struct{
	Slack_name  		string
	Track 				string
	Current_day 		string
	Utc_time 			string
	Github_file_url 	string
	Github_repo_url 	string
	Status_code 		uint
}

type Response struct{
	Error 	bool
	Message string
}

func GetJSONData(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("slack_name")
	track	 := r.URL.Query().Get("track")
	
	if len(username) == 0 || len(track) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")

		var respData Response
		respData.Error = true
		respData.Message = "Slack_name and track cannot be empty"

		resp, err := json.Marshal(respData)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(resp)
		if err != nil {
			panic(err)
		}
		return
	}

	var response responseData

	utc_time := time.Now().UTC().String()

	utcArr := strings.Split(utc_time, " ")
	timeSplit := strings.Split(utcArr[1], ".")
	formattedTime := utcArr[0] + "T" + timeSplit[0] + "Z"


	response.Slack_name = username
	response.Track = track
	response.Current_day = time.Now().Weekday().String()
	response.Utc_time = formattedTime
	response.Github_file_url = "https://github.com/ntekim/hng-stage-1/blob/main/main.go"
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
	mux.Get("/", GetJSONData)

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