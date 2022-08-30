package main

import (
	"encoding/json" // to encode my data into json when sending to Postman
	"fmt"           // formatting package - to print stuff
	"log"           // log out the erros - when errors occur while connecting to server
	"math/rand"     // to create new IDs when user adds a new thing
	"net/http"      // to create a server in go lang
	"strconv"       // to convert IDs to string

	"github.com/gorilla/mux"
)

type Video struct {
	ID      string   `json:"id"`   // video id
	Isbn    string   `json:"isbn"` // unique id assigned to the video
	Title   string   `json:"title"`
	Creator *Creator `json:"creator"`
}

type Creator struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var videos []Video

/*
	//func getVideos( it is response writer; when we send a response 'write' from this function it will be w,

pointer of the request that I will send from my postman to the server)
*/
func getVideos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos) // encodes my go stuff to json
}

func getVideo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range videos {

		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createVideo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var video Video
	_ = json.NewDecoder(r.Body).Decode(&video)

	video.ID = strconv.Itoa(rand.Intn(100000000))
	videos = append(videos, video)
	json.NewEncoder(w).Encode(video)
}

func updateVideo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range videos {

		if item.ID == params["id"] {
			videos = append(videos[:index], videos[index+1:]...)
			var video Video
			_ = json.NewDecoder(r.Body).Decode(&video)

			video.ID = params["id"]
			videos = append(videos, video)
			json.NewEncoder(w).Encode(video)
			return
		}
	}
}

func deleteVideo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // params will have IDs of videos

	for index, item := range videos {

		if item.ID == params["id"] {
			videos = append(videos[:index], videos[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(videos)
}

func main() {
	r := mux.NewRouter()

	videos = append(videos, Video{ID: "1", Isbn: "438227", Title: "Attacks On Titan", Creator: &Creator{Firstname: "Hajime", Lastname: "Isayama"}})
	videos = append(videos, Video{ID: "2", Isbn: "538272", Title: "Grand Blue", Creator: &Creator{Firstname: "Kenji", Lastname: "Kimitake"}})

	//r.HandleFunc("Routes", functionName).Methods("Type?")

	r.HandleFunc("/videos", getVideos).Methods("GET")
	r.HandleFunc("/videos/{id}", getVideo).Methods("GET")
	r.HandleFunc("/videos", createVideo).Methods("POST")
	r.HandleFunc("/videos/{id}", updateVideo).Methods("PUT")
	r.HandleFunc("/videos/{id}", deleteVideo).Methods("DELETE")

	fmt.Println("Starting server at port 8000")

	// log out if function doesn;t start
	log.Fatal(http.ListenAndServe(":8000", r))

}
