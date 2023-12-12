package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Logger *log.Logger
)

func init() {

	fmt.Println("init Logger")

	logDir := "/tmp/"
	logFileName := fmt.Sprintf("RIS_SIMULATION_%s.log", time.Now().UTC().String())
	file, err := os.Create(logDir + logFileName)
	if err != nil {
		log.Panic("could not create log file")
	}

	Logger = log.New(file, "", log.LstdFlags|log.Llongfile)

}
