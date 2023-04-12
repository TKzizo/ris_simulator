package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

const (
	q    float64 = 0.283 // related to the gain
	Gain float64 = math.Pi
	Pt   float64 = 0.05        // Power of transmitter
	P_n  float64 = 0.000000001 // variance of noise at the receiver
)

type Updates struct {
	ris Coordinates
	rx  Coordinates
	tx  Coordinates
	los bool
}

type Simulation struct {
	Env        Environment
	Ris        RIS
	Tx         Tx_Rx
	Rx         Tx_Rx
	Frequency  float64
	f0         float64
	Lambda_p   float64
	Lambda     float64 // wave length
	k          float64
	n_LOS      float64 // Path Loss exponent
	b_LOS      float64 // systemc parameter
	sigma_LOS  float64 // db
	n_NLOS     float64
	b_NLOS     float64
	sigma_NLOS float64
	//channel    chan Updates
	Broadside int8 // 0: SideWall 1: OppositeWall
	Positions []Updates
}

func (s *Simulation) Setup() {
	s.Lambda = 3.0 / 10 * s.Frequency // it's Simplified so it only supports GHz
	s.k = 2 * math.Pi / s.Lambda

	if s.Frequency == 28.0 {
		s.Lambda_p = 1.8
	} else if s.Frequency == 73.0 {
		s.Lambda_p = 1.9
	}

	if s.f0 == 0.0 {
		s.f0 = 24.2
	}

	if s.n_LOS == 0.0 { //Pathloss exponent
		s.n_LOS = 1.73
	}

	if s.n_NLOS == 0.0 {
		s.n_NLOS = 3.79
	}

	if s.b_NLOS == 0.0 {
		s.b_NLOS = 3.19
	}

	if s.sigma_LOS == 0.0 {
		s.sigma_LOS = 3.02
	}

	if s.sigma_NLOS == 0.0 {
		s.sigma_NLOS = 8.29
	}

	s.Ris.Setup(s.Lambda)
	s.InputPositions()
	//s.CheckPositioning()
}

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
	for _, v := range list_positions {
		fmt.Println(v)
	}
	s.Positions = list_positions

}

func (s *Simulation) Rate_SNR(H, G, Theta mat.CDense, SNR int) (float64, float64) {

	var temp1 mat.CDense
	var temp2 mat.CDense

	temp1.Mul(G.T(), &Theta)
	temp2.Mul(&temp1, &H)
	fmt.Println(temp2.RawCMatrix().Data)

	snr := math.Pow(cmplx.Abs(temp2.RawCMatrix().Data[0]), 2) * Pt / P_n

	return math.Log2(1 + snr), float64(SNR) * snr
}
