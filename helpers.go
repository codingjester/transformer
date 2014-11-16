package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func trimExtension(filename string) string {
	extension := filepath.Ext(filename)
	return strings.TrimSuffix(filename, extension)
}

func moveToJBsMagicHome(srcFileName string, outputFileName string) {
	filename := fmt.Sprintf("/Users/johnb/Desktop/Transformer/%s", outputFileName)
	os.Rename(srcFileName, filename)
}

func insert_job_id(job_id string) {
	stmt, err := db.Prepare("INSERT INTO jobs (job_id) VALUES(?)")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = stmt.Exec(job_id)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func start_job_id(job_id string) {
	stmt, err := db.Prepare("UPDATE jobs SET status = 2 WHERE job_id = ?")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = stmt.Exec(job_id)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println(fmt.Sprintf("Starting job_id %s", job_id))
}

func finish_job_id(job_id string) {
	stmt, err := db.Prepare("UPDATE jobs SET status = 3 WHERE job_id = ?")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = stmt.Exec(job_id)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println(fmt.Sprintf("Finishing job_id %s", job_id))
}

func generate_job_id(length int) string {
	rb := make([]byte, length)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}
	rs := base64.URLEncoding.EncodeToString(rb)
	return rs[:length]
}

func WriteJSON(w http.ResponseWriter, js []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
