package main

import "fmt"

func main() {

	ris := RIS{N: 256, xyz: Coordinates{x: 40, y: 50, z: 2}}
	tx := Tx_Rx{N: 1, xyz: Coordinates{x: 0, y: 25, z: 2}}
	rx := Tx_Rx{N: 1, xyz: Coordinates{x: 38, y: 48, z: 1}}

	simulation := Simulation{Ris: ris, Tx: tx, Rx: rx, Frequency: 28.0, channel: make(chan Updates, 1)}

	simulation.Setup()
	h := simulation.H_channel()
	g := simulation.G_channel()
	theta := GetCoefficients(h, g)

	rate := simulation.Rate(h, g, theta)

	fmt.Println("Rate: ", rate)
}
