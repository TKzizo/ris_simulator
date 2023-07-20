package main

import (
	"encoding/json"
	"io/ioutil"
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
	Tx  tx_rx_ris `json:"TX"`
	Rx  tx_rx_ris `json:"RX"`
	Ris tx_rx_ris `json:"RIS"`
}
type tx_rx_ris struct {
	Elements int `json:"Elements"`
	Coord    struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"Coordinates"`
	Type int `json:"Type"`
}

func initCfg(cfgPath string) InitConfig {

	var cfg InitConfig

	jsonfile, err := os.Open(cfgPath)
	if err != nil {
		log.Println(err)
	}
	bytes, err := ioutil.ReadAll(jsonfile)
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(bytes, &cfg)

	defer jsonfile.Close()

	return cfg

}
