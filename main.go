package main

import (
	"fmt"
	"os"
)

func main() {

	ris := RIS{N: 16, xyz: Coordinates{x: 40, y: 50, z: 2}}
	tx := Tx_Rx{N: 2, xyz: Coordinates{x: 0, y: 25, z: 2}}
	rx := Tx_Rx{N: 2, xyz: Coordinates{x: 38, y: 48, z: 1}}

	simulation := Simulation{Ris: ris, Tx: tx, Rx: rx, Frequency: 28.0, Env: Environment{75.0, 50.0, 3.5}}

	simulation.Setup()
	list := simulation.Run()

	for _, mat := range list {
		//fmt.Println(mat)
		fmt.Fprintln(os.Stderr, mat)
	}

}
