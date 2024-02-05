package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type InitConfig struct {
	Frequency float64 `json:"Frequency"`
	Env       struct {
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"Environment"`
	Equipements _Equipements `json:"Equipements"`
}

type _Equipements struct {
	Ris ris   `json:"RIS"`
	Tx  tx_rx `json:"TX"`
	Rx  tx_rx `json:"RX"`
}

type tx_rx struct {
	Elements int `json:"NbrElements"`
	Coord    struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"Position"`
	Type int `json:"Type"`
}

type ris struct {
	Elements int `json:"NbrElements"`
	Coord    struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"Position"`
	Broadside int `json:"Broadside"`
}

func InitCfg(cfgPath string) InitConfig {

	var cfg InitConfig

	jsonfile, err := os.Open(cfgPath)
	if err != nil {
		log.Println(err)
	}
	bytes, err := io.ReadAll(jsonfile)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(bytes, &cfg)

	defer jsonfile.Close()

	return cfg

}
