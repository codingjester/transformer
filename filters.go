package main

import (
	"log"
	"os/exec"
)

func Lomo(srcFileName string, outputFileName string) {

	cmd := exec.Command("convert", srcFileName, "-channel", "R", "-level", "33%", "-channel", "G", "-level", "33%", srcFileName)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Waiting for the command to finish")
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done.")

}

func Gotham(srcFileName string, outputFileName string) {

	cmd := exec.Command("convert", srcFileName, "-modulate", "120,10,100", "-fill", "#222b6d", "-colorize", "20", "-gamma", "0.5", "-contrast", "-contrast", srcFileName)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Waiting for the command to finish")
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done.")

}

func Toaster(srcFileName string, outputFileName string) {
	cmd := exec.Command("convert", srcFileName, "-modulate", "150,80,100", "-gamma", "1.2", "-contrast", "-contrast", srcFileName)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Waiting for the command to finish")
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done.")
}
