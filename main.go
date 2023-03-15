package main

import (
	"fmt"
)

func main() {

	/*c := mat.NewDense(4, 1, []float64{1, 2, 3, 4})
	fc := mat.Formatted(c, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("with all values:\na = %v\n\n", fc)*/
	ris := RIS{N: 64, xyz: Coordinates{x: 40, y: 50, z: 2}}
	tx := Tx_Rx{N: 1, xyz: Coordinates{x: 0, y: 25, z: 2}}
	rx := Tx_Rx{N: 1, xyz: Coordinates{x: 38, y: 48, z: 1}}

	simulation := Simulation{Ris: ris, Tx: tx, Rx: rx, Frequency: 28.0, channel: make(chan Updates, 1)}

	simulation.Setup()

	h := simulation.H_channel()
	g := simulation.G_channel()
	theta := GetCoefficients(h, g)
	fmt.Println(h)
	fmt.Println(g)
	fmt.Println(theta)

	rate, snr := simulation.Rate_SNR(h, g, theta, 1)
	fmt.Println(rate, snr)
}
