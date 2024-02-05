package simulate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"log"
)

func PostionHandler(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(simulation.Positions)
	w.Write(j)
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(cfg)
	w.Write(j)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("web/index.html")
	if err != nil {
		log.Print(err)
	}

	fmt.Println(string(file))
	w.Write(file)
}
