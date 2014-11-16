package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func ApplyFilter(job_id string, srcFileName string, outputFileName string, filter string) {
	insert_job_id(job_id)
	start_job_id(job_id)
	if filter == "gotham" {
		Gotham(srcFileName, outputFileName)
	} else if filter == "lomo" {
		Lomo(srcFileName, outputFileName)
	} else if filter == "toaster" {
		Toaster(srcFileName, outputFileName)
	}
	moveToJBsMagicHome(srcFileName, outputFileName)
	finish_job_id(job_id)
	go FireDehMissels(job_id, outputFileName)
}

func TranscodeGif(job_id string, srcFileName string, outputFileName string) {
	log.Println(job_id)

	insert_job_id(job_id)

	file_path := fmt.Sprintf("/tmp/%s", job_id)
	os.Mkdir(file_path, 0777)
	outfile := fmt.Sprintf("%s/%s", file_path, "tmp.mp4")
	cmd := exec.Command("ffmpeg", "-i", srcFileName, "-pix_fmt", "yuv420p", "-vf", "scale=trunc(in_w/2)*2:trunc(in_h/2)*2'", "-r", "30", outfile)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Waiting for the command to finish")
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	moveToJBsMagicHome(outfile, outputFileName)
	finish_job_id(job_id)
	go FireDehMissels(job_id, outputFileName)
}

func Transcode(job_id string, srcFileName string, outputFileName string) {

	insert_job_id(job_id)

	start_job_id(job_id)
	// Do the actual ffmpeg transcode for 720 and 480
	doTranscoding(srcFileName, outputFileName, "hd720")
	//doTranscoding(srcFileName, outputFileName, "hd480")

	// cleanup the src filename
	os.Remove(srcFileName)
	finish_job_id(job_id)
	// Run cleanup
	log.Printf("Done.")
	go FireDehMissels(job_id, outputFileName)

}

func doTranscoding(srcFileName string, outputFileName string, sizing string) {

	log.Println(fmt.Sprintf("Doing the %s", sizing))
	outputfile := fmt.Sprintf("%s.mp4", srcFileName)
	cmd := exec.Command("ffmpeg", "-i", srcFileName, "-s", sizing, "-c:v", "libx264", "-crf", "23", "-c:a", "aac", "-strict", "-2", "-f", "mp4", outputfile)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Waiting for the command to finish")
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	ExtractFrames(outputfile, outputFileName)
	moveToJBsMagicHome(outputfile, outputFileName)

}

func ExtractFrames(srcFileName string, outputFileName string) {
	log.Printf("Extracting thumbnails")
	thumbnail_name := trimExtension(outputFileName)
	outputdir := fmt.Sprintf("/Users/johnb/Desktop/Transformer/%s_thumb%%03d.jpg", thumbnail_name)
	log.Printf(outputdir)
	cmd := exec.Command("ffmpeg", "-i", srcFileName, "-r", "1/10", outputdir)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Waiting for the command to finish")
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	file := fmt.Sprintf("/Users/johnb/Desktop/Transformer/%s_thumb*.jpg", thumbnail_name)
	filmstripfile := fmt.Sprintf("/Users/johnb/Desktop/Transformer/%s_filmstrip.jpg", thumbnail_name)

	cmd = exec.Command("convert", file, "+append", "-quality", "70", filmstripfile)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Waiting for the filmstrip to finish")
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done extracting thumbnails")

}

type DehMissel struct {
	Job_Id string `json:"job_id"`
	File   string `json:"file"`
}

func FireDehMissels(job_id string, filename string) {

	missel := DehMissel{job_id, filename}
	json, err := json.Marshal(missel)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:4567/ping", bytes.NewBuffer(json))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("All dah missels fired")
}
