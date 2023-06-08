package main

import (
	cmat "RIS_SIMULATOR/reducedComplex"
	"math"
	"math/cmplx"

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
	s.Lambda = 3.0 / (10 * s.Frequency) // it's Simplified so it only supports GHz
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
	s.Rx.Setup(s.Lambda)
	s.Tx.Setup(s.Lambda)
	s.InputPositions()
	//s.CheckPositioning() // To apply the 3GPP standards
}

func (s *Simulation) rate(H, G mat.CDense, Theta mat.CDiagonal) float64 {

	var temp1 mat.CDense
	var temp2 mat.CDense
	rate := 0.0

	temp1.Mul(G.T(), Theta)
	temp2.Mul(&temp1, &H)
	rate = math.Log2(math.Pow(cmplx.Abs(temp2.At(0, 0)), 2) * Pt / P_n)

	return rate
}

func (s *Simulation) Run() (*cmat.Cmatrix, *cmat.Cmatrix) {
	var h, g cmat.Cmatrix
	//list := []cmat.Cmatrix{h, g}
	// Re-run the calculation for every position of the user
	for _, update := range s.Positions {
		clusters := GenerateClusters(s)
		s.Ris.xyz = update.ris
		s.Tx.xyz = update.tx
		s.Rx.xyz = update.rx
		h = s.H_channel(clusters)
		//list = append(list, h)
		g = s.G_channel()
		//list = append(list, g)
	}
	return &h, &g

}

/*func (s *Simulation) MIMO_Rate(H, G mat.CDense, Theta mat.CDiagonal) float64 {
	var temp1 mat.CDense
	var temp2 mat.CDense
	rate := 0.0

	temp1.Mul(G.T(), Theta)
	temp2.Mul(&temp1, &H)
	temp1.Mul(temp2.H(), &temp2)
	for i := 0; i < temp1.RawCMatrix().Rows; i++ {
		temp1.Set(i, i, temp1.At(i, i)+complex(1, 0))
	}
	temp1.Scale(complex(Pt/P_n, 0), &temp1)
	var lu cLU
	lu.Factorize(temp1)
	rate = lu.Det()

	return rate

}*/
