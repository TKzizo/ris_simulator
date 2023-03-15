package main

type Tx_Rx struct {
	N int //Number of antenna elements
	//Type int (1 for planar and 0 for linear)
	xyz       Coordinates
	Theta_RIS float64
	Theta     float64
	Phi_RIS   float64
	Phi       float64
}
