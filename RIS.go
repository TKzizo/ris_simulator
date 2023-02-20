package main

type RIS struct {
	N        int
	dis      float64 // distance between elements
	xyz      Coordinates
	Theta_Tx float64
	Theta_Rx float64
	Phi_Tx   float64
	Phi_Rx   float64
}

func (r *RIS) Setup(Lambda float64) {
	if r.N == 0 {
		r.N = 256
	}
	if r.dis == 0.0 {
		r.dis = Lambda / 2
	}
}
