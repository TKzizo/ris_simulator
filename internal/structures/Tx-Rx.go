package structures

type Tx_Rx struct {
	N         int     //Number of antenna elements
	Type      int     // (1 for planar and 0 for linear)
	Dis       float64 // distance between antennas
	Xyz       Coordinates
	Theta_RIS float64
	Phi_RIS   float64
	// elevation and azimith Rx->Tx and Tx->RX
	Theta float64
	Phi   float64
}

func (r *Tx_Rx) Setup(Lambda float64) {
	if r.Dis == 0.0 {
		r.Dis = Lambda / 2
	}
}
