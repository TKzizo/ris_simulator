package utils

import (
	"log"
	"os"
)

func checkDir(directory string) error {
	_, err := os.Stat(directory)
	return err
}

func CreateDir(directory string) error {
	log.Println("Creating Directory: ", directory)
	err := os.Mkdir(directory, os.ModePerm)
	if err != nil {
		log.Println("Could not create Directory: ", directory)
		return err
	}
	log.Println("Directory created successfully")
	return nil
}

func DefaultDir(directory string) error {
	err := checkDir(directory)
	if os.IsNotExist(err) {
		err = CreateDir(directory)
		return err
	}
	return err
}
