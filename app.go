package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"runtime"
	"runtime/debug"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func assert(i interface{}, err error) interface{} {
	if err != nil {
		fatal(err)
	}
	return i
}

func fatal(err error) {
	debug.PrintStack()
	log.Fatal(err)
	os.Exit(0)
}

type testResponse struct {
	Response string `json:"job_id"`
}

type StatusResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

var db *sql.DB
var config *Configuration

type Configuration struct {
	Hostname    string
	Proto       string
	Port        int
	Db_Type     string
	Db_Username string
	Db_Password string
	Db_Host     string
	DB          string
}

func AcceptConvert(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	outputfile := r.FormValue("outputfile")
	filter := r.FormValue("filter")

	log.Println(outputfile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp, err := ioutil.TempFile(os.TempDir(), "image_")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tmp, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file.Close()
	tmp.Close()

	job_id := generate_job_id(50)
	go ApplyFilter(job_id, tmp.Name(), outputfile, filter)

	resp := testResponse{job_id}
	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, json)
}

func AcceptTranscode(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("video_file")
	outputfile := r.FormValue("outputfile")

	log.Println(outputfile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp, err := ioutil.TempFile(os.TempDir(), "video_")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tmp, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file.Close()
	tmp.Close()

	job_id := generate_job_id(50)
	go Transcode(job_id, tmp.Name(), outputfile)

	resp := testResponse{job_id}
	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, json)
}

func AcceptGifTranscode(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("image")
	outputfile := r.FormValue("outputfile")

	log.Println(outputfile)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp, err := ioutil.TempFile(os.TempDir(), "image_")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tmp, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file.Close()
	tmp.Close()

	job_id := generate_job_id(50)
	go TranscodeGif(job_id, tmp.Name(), outputfile)

	resp := testResponse{job_id}
	json, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, json)
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hash := params["hash"]

	var status_id int
	err := db.QueryRow("SELECT status FROM jobs WHERE job_id = ?", hash).Scan(&status_id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var status string
	if status_id == 1 {
		status = "Starting"
	} else if status_id == 2 {
		status = "Running"
	} else if status_id == 3 {
		status = "Completed"
	} else {
		status = "Unknown"
	}
	resp := StatusResponse{hash, status}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, js)
}

func loadConfig() {
	file, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatal("unable to open config: ", err)
	}

	temp := new(Configuration) // Get a pointer to an instance with new keyword
	// Unmarshal is going to decode and store into temp
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config", err)
	}
	config = temp
}

// Sets up the datbase using the loaded config
// Possible improvements would be adding in other database support
func setupDB() {

	var err error
	database_url := fmt.Sprintf("%s@%s/%s", config.Db_Username, config.Db_Host, config.DB)
	db, err = sql.Open(config.Db_Type, database_url)
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}

	db.SetMaxIdleConns(10)
	err = db.Ping() // Check for DB access
	if err != nil {
		log.Fatalf("Error on opening database connection: %s", err.Error())
	}

}

func main() {
	runtime.GOMAXPROCS(4)

	loadConfig()

	setupDB()

	r := mux.NewRouter()
	r.HandleFunc("/transcode", AcceptTranscode).Methods("POST")
	r.HandleFunc("/convert", AcceptConvert).Methods("POST")
	r.HandleFunc("/gif", AcceptGifTranscode).Methods("POST")
	r.HandleFunc("/status/{hash}", GetStatus).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
