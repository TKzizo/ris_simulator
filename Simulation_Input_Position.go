package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func (s *Simulation) InputPositions() {
	csvFile, err := os.Open("Positions.csv")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("number of lines", len(csvLines))
	list_positions := []Updates{}

	for _, line := range csvLines {
		position := Updates{}
		if v, err := strconv.ParseFloat(line[0], 64); err == nil && v <= s.Env.length && v >= 0 {
			position.ris.x = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}
		if v, err := strconv.ParseFloat(line[1], 64); err == nil && v <= s.Env.width && v >= 0 {
			position.ris.y = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[2], 64); err == nil && v <= s.Env.height && v >= 0 {
			position.ris.z = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[3], 64); err == nil && v <= s.Env.length && v >= 0 {
			position.tx.x = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[4], 64); err == nil && v <= s.Env.width && v >= 0 {
			position.tx.y = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[5], 64); err == nil && v <= s.Env.height && v >= 0 {
			position.tx.z = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[6], 64); err == nil && v <= s.Env.length && v >= 0 {
			position.rx.x = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[7], 64); err == nil && v <= s.Env.width && v >= 0 {
			position.rx.y = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[8], 64); err == nil && v <= s.Env.height && v >= 0 {
			position.rx.z = v
		} else {
			fmt.Println("Position Over-Boundries")
			continue
		}

		if v, err := strconv.ParseFloat(line[9], 64); err == nil {
			if v == 1 {
				position.los = true
			} else {
				position.los = false
			}
		}
		list_positions = append(list_positions, position)
	}
	/*for _, v := range list_positions {
		fmt.Println(v)
	}*/
	s.Positions = list_positions

}
