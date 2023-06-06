package main

type RIS struct {
	N        int
	dis      float64 // distance between elements
	xyz      Coordinates
	Theta_Tx float64
	Theta_Rx float64
	Phi_Tx   float64
	Phi_Rx   float64
	//nbits  int // number of bits -> number of reflection states
}

func (r *RIS) Setup(Lambda float64) {

	if r.N == 0 { // TO be removed since sim must precise the number of patches
		r.N = 256
	}
	if r.dis == 0.0 {
		r.dis = Lambda / 2
	}
}
