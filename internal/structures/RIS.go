package structures

type RIS struct {
	N        int
	Dis      float64 // distance between elements
	Xyz      Coordinates
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
	if r.Dis == 0.0 {
		r.Dis = Lambda / 2
	}
}
